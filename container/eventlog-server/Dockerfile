From golang:1.20-alpine3.17 AS builder

RUN apk update \
    && apk add --no-cache protoc make wget

WORKDIR /usr/bin
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.18 && \
    wget -qO grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x grpc_health_probe

WORKDIR /go/src/github.com/eventlog-server
COPY service/eventlog-server ./
RUN make all

From alpine:3.17.0

ARG USER=ccnp
ARG UID=1000
ARG GID=1000
ARG GROUP=ccnp

WORKDIR /opt/eventlog-server
RUN addgroup -S -g $GID $GROUP && adduser -S -u $UID -D -G $GROUP $USER
RUN chown $USER:$GROUP /opt/eventlog-server

COPY --chown=$USER --from=builder /go/src/github.com/eventlog-server/eventlog-server ./
COPY --chown=$USER --from=builder /usr/bin/grpc_health_probe /usr/bin/grpc_health_probe

USER $UID

CMD ["./eventlog-server"]
