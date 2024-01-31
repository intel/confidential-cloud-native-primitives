import os

import grpc


import ccnp_server_pb2_grpc
import ccnp_server_pb2

DEFAULT_SOCK: str = "unix:/run/ccnp/uds/ccnp-server.sock"

class CCNPClient:
    def __init__(self, sock: str = DEFAULT_SOCK):
        self._sock = sock
        
    def GetReport(self, user_data: str, nonce: str) -> ccnp_server_pb2.GetReportResponse:
        stub = self._get_stub()
        req = ccnp_server_pb2.GetReportRequest(user_data=user_data, nonce=nonce)
        return stub.GetReport(req)
        
    def GetMeasurement(self, index: int, algo_id: int) -> ccnp_server_pb2.GetMeasurementResponse:
        stub = self._get_stub()
        req = ccnp_server_pb2.GetMeasurementRequest(index=index, algo_id=algo_id)
        return stub.GetMeasurement(req)
        
    def GetEventlog(self, start: int, count: int) -> ccnp_server_pb2.GetEventlogResponse:
        stub = self._get_stub()
        req = ccnp_server_pb2.GetEventlogRequest(start=start, count=count)
        return stub.GetEventlog(req)
        

    def _get_stub(self) -> ccnp_server_pb2_grpc.ccnpStub:
        if not os.path.exists(self._sock.replace('unix:', '')):
            raise RuntimeError("CCNP server does not start.")
        channel = grpc.insecure_channel(self._sock,
                                        options=[('grpc.default_authority', 'localhost')])
        stub = ccnp_server_pb2_grpc.ccnpStub(channel)
        return stub


# python3 -m grpc_tools.protoc -I proto --python_out=py_gen/ --pyi_out=py_gen/ --grpc_python_out=py_gen/ proto/ccnp-server.proto
