# CCNP service


## Query info

1. Query the Report

Run the command.

```
grpcurl -authority "dummy"  -plaintext -d '{"level": 1, "user_data": "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4", "nonce":"IXUKoBO1UM3c1wopN4sY" }'  -unix /run/ccnp/uds/ccnp-server.sock ccnp_server_pb.ccnp.GetReport
```

The output looks like this.

```
{
    "report": "..."
}
```

2. Query the measurement

Run the command.

```
grpcurl -authority "dummy"  -plaintext -d '{"level": 1, "index": 0}'  -unix /run/ccnp/uds/ccnp-server.sock ccnp_server_pb.ccnp.GetMeasurement
```

The output looks like.

```
{
  "measurement": "..."
}
```

3. Query the eventlog

Run the command.

```
grpcurl -authority "dummy"  -plaintext -d '{"level": 1, "start": 1, "count": 3}'  -unix /run/ccnp/uds/ccnp-server.sock ccnp_server_pb.ccnp.GetEventlog
```

The output looks like.

```
{
  "events": [
    {
      "eventType": 3,
      "digest": "AAAAAAAAAAAAAAAAAAAAAAAAAAA=",
      "eventSize": 33,
      "event": "U3BlYyBJRCBFdmVudDAzAAAAAAAAAgACAQAAAAwAMAAA"
    },
    {
      "eventType": 2147483659,
      "eventSize": 42,
      "event": "CVRkeFRhYmxlAAEAAAAAAAAAr5a7k/K5uE6UYuC6dFZCNgCQgAAAAAAA",
      "digests": [
        {
          "algoId": 12,
          "hash": "C4dy5bC0G4PmBEpoOX4C9J+0cGa0++SRfqLEXGTzI/2suzeUj4Ieuvi8nJOLqKdJ"
        }
      ]
    },
    {
      "eventType": 2147483658,
      "eventSize": 58,
      "event": "KUZ2KFhYWFhYWFhYLVhYWFgtWFhYWC1YWFhYLVhYWFhYWFhYWFhYWCkAAADA/wAAAAAAQAgAAAAAAA==",
      "digests": [
        {
          "algoId": 12,
          "hash": "NEvFHJgLpiGqoA2j7XQ299blSRl9/mmVFd+ixlg9leZBKvIcCX1HMVWHX/1WHWeQ"
        }
      ]
    }
  ]
}
```

## Server log

The server will print all requests and responses and assigns an uuid to the pair of request and response.

It looks like this.

```
./ccnp_server
Starting ccnp server in "TDX" enviroment...
2024-01-24T03:21:07.130Z INFO  [ccnp_server::handler] Request IN ---> id:42e1ff55-48bf-41e1-bf05-b5e67020a251: GetReportRequest { level: Paas, user_data: "MT}
2024-01-24T03:21:07.131Z INFO  [cctrusted_vm::tdvm] ======================================
2024-01-24T03:21:07.131Z INFO  [cctrusted_vm::tdvm] CVM type = TDX
2024-01-24T03:21:07.131Z INFO  [cctrusted_vm::tdvm] CVM version = 1.5
2024-01-24T03:21:07.131Z INFO  [cctrusted_vm::tdvm] ======================================
2024-01-24T03:21:07.139Z INFO  [ccnp_server::handler] Response OK <--- id:42e1ff55-48bf-41e1-bf05-b5e67020a251
2024-01-24T03:21:14.635Z INFO  [ccnp_server::handler] Request IN ---> id:3bcdbbfa-9d53-4eb2-a145-33cc8661756b: GetMeasurementRequest { level: Paas, index: 0 }
2024-01-24T03:21:14.636Z INFO  [ccnp_server::handler] Response OK <--- id:3bcdbbfa-9d53-4eb2-a145-33cc8661756b
2024-01-24T03:21:20.822Z INFO  [ccnp_server::handler] Request IN ---> id:d2894a22-2c89-4016-82e4-c6dfea8bafcf: GetEventlogRequest { level: Paas, start: 1, co}
2024-01-24T03:21:20.823Z INFO  [ccnp_server::handler] Response OK <--- id:d2894a22-2c89-4016-82e4-c6dfea8bafcf

```
