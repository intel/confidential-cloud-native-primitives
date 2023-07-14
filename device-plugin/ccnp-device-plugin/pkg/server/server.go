/* SPDX-license-identifier: Apache-2.0 */
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	dpapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	Namespace                  = "tdx.intel.com"
	DeviceType                 = "tdx-guest"
	CcnpDpSocket               = "/var/lib/kubelet/device-plugins/ccnpdp.sock"
	KubeletSocket              = "/var/lib/kubelet/device-plugins/kubelet.sock"
	TDX_DEVICE_DEPRECATED      = "/dev/tdx-attest"
	TDX_DEVICE_1_0             = "/dev/tdx-guest"
	TDX_DEVICE_1_5             = "/dev/tdx_guest"
	TdxDevicePermissions       = "rw"
	MaxRestartCount            = 5
	SocketConnectTimeout       = 5
	DefaultPodCount       uint = 110
	UDS_WORK_DIR               = "/run/ccnp/uds"
)

type CcnpDpServer struct {
	srv            *grpc.Server
	devices        map[string]*dpapi.Device
	ctx            context.Context
	cancel         context.CancelFunc
	restartFlag    bool
	tdxGuestDevice string
}

func NewCcnpDpServer() *CcnpDpServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &CcnpDpServer{
		devices:     make(map[string]*dpapi.Device),
		srv:         grpc.NewServer(grpc.EmptyServerOption{}),
		ctx:         ctx,
		cancel:      cancel,
		restartFlag: false,
	}
}

func (ccnpdpsrv *CcnpDpServer) getTdxVersion() error {

	if _, err := os.Stat(TDX_DEVICE_DEPRECATED); err == nil {
		return errors.New("Deprecated TDX device found")
	}

	if _, err := os.Stat(TDX_DEVICE_1_0); err == nil {
		ccnpdpsrv.tdxGuestDevice = TDX_DEVICE_1_0
		return nil
	}

	if _, err := os.Stat(TDX_DEVICE_1_5); err == nil {
		ccnpdpsrv.tdxGuestDevice = TDX_DEVICE_1_5
		return nil
	}

	return errors.New("No TDX device found")
}

func (ccnpdpsrv *CcnpDpServer) scanDevice() error {

	err := ccnpdpsrv.getTdxVersion()
	if err != nil {
		return err
	}

	for i := uint(0); i < DefaultPodCount; i++ {
		deviceID := fmt.Sprintf("%s-%d", "tdx-guest", i)
		ccnpdpsrv.devices[deviceID] = &dpapi.Device{
			ID:     deviceID,
			Health: dpapi.Healthy,
		}
	}

	return nil
}

func (ccnpdpsrv *CcnpDpServer) Run() error {

	err := ccnpdpsrv.scanDevice()
	if err != nil {
		klog.Fatalf("scan device error: %v", err)
	}

	dpapi.RegisterDevicePluginServer(ccnpdpsrv.srv, ccnpdpsrv)

	err = syscall.Unlink(CcnpDpSocket)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	listen, err := net.Listen("unix", CcnpDpSocket)
	if err != nil {
		return err
	}

	go func() {
		failCount := 0
		for {
			err = ccnpdpsrv.srv.Serve(listen)
			if err == nil {
				break
			}

			if failCount > MaxRestartCount {
				klog.Fatalf("CCNP plugin server crashed. Quitting...")
			}
			failCount++
		}
	}()

	connection, err := ccnpdpsrv.connect(CcnpDpSocket, time.Duration(SocketConnectTimeout)*time.Second)
	if err != nil {
		return err
	}

	connection.Close()

	return nil
}

func (s *CcnpDpServer) connect(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {

	connection, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (ccnpdpsrv *CcnpDpServer) RegisterToKubelet() error {

	conn, err := ccnpdpsrv.connect(KubeletSocket, time.Duration(MaxRestartCount)*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := dpapi.NewRegistrationClient(conn)
	request := &dpapi.RegisterRequest{
		Version:      dpapi.Version,
		Endpoint:     path.Base(CcnpDpSocket),
		ResourceName: Namespace + "/" + DeviceType,
	}

	_, err = client.Register(context.Background(), request)
	if err != nil {
		return err
	}

	return nil
}

func (ccnpdpsrv *CcnpDpServer) ListAndWatch(e *dpapi.Empty, lwSrv dpapi.DevicePlugin_ListAndWatchServer) error {
	tdxDevices := make([]*dpapi.Device, len(ccnpdpsrv.devices))

	i := 0
	for _, tdxDevice := range ccnpdpsrv.devices {
		tdxDevices[i] = tdxDevice
		i++
	}

	err := lwSrv.Send(&dpapi.ListAndWatchResponse{Devices: tdxDevices})
	if err != nil {
		klog.Fatalf("ListAndWatch error: %v", err)
		return err
	}

	for {
		select {
		case <-ccnpdpsrv.ctx.Done():
			return nil
		}
	}
}

func (ccnpdpsrv *CcnpDpServer) GetDevicePluginOptions(ctx context.Context, e *dpapi.Empty) (*dpapi.DevicePluginOptions, error) {
	return &dpapi.DevicePluginOptions{PreStartRequired: true}, nil
}

func (ccnpdpsrv *CcnpDpServer) GetPreferredAllocation(ctx context.Context, r *dpapi.PreferredAllocationRequest) (*dpapi.PreferredAllocationResponse, error) {
	return &dpapi.PreferredAllocationResponse{}, nil
}

func (ccnpdpsrv *CcnpDpServer) PreStartContainer(ctx context.Context, req *dpapi.PreStartContainerRequest) (*dpapi.PreStartContainerResponse, error) {
	return &dpapi.PreStartContainerResponse{}, nil
}

func (ccnpdpsrv *CcnpDpServer) Allocate(ctx context.Context, reqs *dpapi.AllocateRequest) (*dpapi.AllocateResponse, error) {
	response := &dpapi.AllocateResponse{}

	devSpec := dpapi.DeviceSpec{
		HostPath:      ccnpdpsrv.tdxGuestDevice,
		ContainerPath: ccnpdpsrv.tdxGuestDevice,
		Permissions:   TdxDevicePermissions,
	}

	pluginMount := dpapi.Mount{
		ContainerPath: UDS_WORK_DIR,
		HostPath:      UDS_WORK_DIR,
	}

	for range reqs.ContainerRequests {
		klog.Infof("received resource request")
		resp := dpapi.ContainerAllocateResponse{
			Envs:        make(map[string]string),
			Annotations: make(map[string]string),
			Devices:     []*dpapi.DeviceSpec{&devSpec},
			Mounts:      []*dpapi.Mount{&pluginMount},
		}
		response.ContainerResponses = append(response.ContainerResponses, &resp)
	}
	return response, nil
}
