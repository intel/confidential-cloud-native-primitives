/* SPDX-license-identifier: Apache-2.0 */
package main

import (
	"log"
	"os"
	"path"

	"ccnp-device-plugin/pkg/server"

	"github.com/fsnotify/fsnotify"
	"k8s.io/klog/v2"
)

func main() {

	log.Println("Intel CCNP device plugin starting")
	ccnpdpsrv := server.NewCcnpDpServer()
	go ccnpdpsrv.Run()

	if err := ccnpdpsrv.RegisterToKubelet(); err != nil {
		klog.Errorf("register to kubelet error: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Fatalf("Failed to created FS watcher.")
		os.Exit(1)
	}
	defer watcher.Close()

	err = watcher.Add(path.Dir(server.KubeletSocket))
	if err != nil {
		klog.Fatalf("watch kubelet error")
		return
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Name == server.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				klog.Fatalf("restart CCNP device plugin due to kubelet restart")
			}
			if event.Name == server.CcnpDpSocket && event.Op&fsnotify.Remove == fsnotify.Remove {
				klog.Fatalf("restart CCNP device plugin due to device plugin socket being deleted")
			}
		case err := <-watcher.Errors:
			klog.Fatalf("fsnotify watch error: %s", err)
		}
	}
}
