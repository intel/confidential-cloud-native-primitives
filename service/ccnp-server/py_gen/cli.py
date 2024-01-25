import os
import ccnp_server_pb2_grpc
import ccnp_server_pb2

import grpc

DEFAULT_SOCK: str = "unix:/run/ccnp/uds/ccnp-server.sock"

class CCNPClient:
    def __init__(self, sock: str = DEFAULT_SOCK):
        self._sock = sock
        
    def GetReport(self, level: ccnp_server_pb2.LEVEL, user_data: str, nonce: str) -> ccnp_server_pb2.GetReportResponse:
        stub = self._get_stub()
        req = ccnp_server_pb2.GetReportRequest(level=level, user_data=user_data, nonce=nonce)
        return stub.GetReport(req)
        
    def GetMeasurement(self, level: ccnp_server_pb2.LEVEL, index: int) -> ccnp_server_pb2.GetMeasurementResponse:
        stub = self._get_stub()
        req = ccnp_server_pb2.GetMeasurementRequest(level=level, index=index)
        return stub.GetMeasurement(req)
        
    def GetEventlog(self, level: ccnp_server_pb2.LEVEL, start: int, count: int) -> ccnp_server_pb2.GetEventlogResponse:
        stub = self._get_stub()
        req = ccnp_server_pb2.GetEventlogRequest(level=level, start=start, count=count)
        return stub.GetEventlog(req)
        

    def _get_stub(self) -> ccnp_server_pb2_grpc.ccnpStub:
        if not os.path.exists(self._sock.replace('unix:', '')):
            raise RuntimeError("CCNP server does not start.")
        channel = grpc.insecure_channel(self._sock,
                                        options=[('grpc.default_authority', 'localhost')])
        stub = ccnp_server_pb2_grpc.ccnpStub(channel)
        return stub

if __name__ == "__main__":
    cli = CCNPClient()
    # resp = cli.GetReport(ccnp_server_pb2.LEVEL.PAAS, "", "")
    # print(resp.report)
    # resp = cli.GetMeasurement(ccnp_server_pb2.LEVEL.PAAS, 0)
    # print(resp.measurement)
    resp = cli.GetEventlog(ccnp_server_pb2.LEVEL.PAAS, 1, 3)
    print(resp.events)


# python3 -m grpc_tools.protoc -I proto --python_out=py_gen/ --pyi_out=py_gen/ --grpc_python_out=py_gen/ proto/ccnp-server.proto
