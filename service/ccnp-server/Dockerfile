FROM rust:1.71-alpine3.18 AS quote-server-builder

RUN apk update \
    && apk add --no-cache make protobuf-dev musl-dev wget

COPY . /ccnp-server
RUN cd /ccnp-server && make build

# add rediness and liveness probe command
WORKDIR /usr/bin
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.22 && \
    wget -qO grpc-health-probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x grpc-health-probe

# ======================================================================================================================

FROM alpine:3.18.5
ARG USER=ccnp
ARG GROUP=ccnp
ARG UID=1000
ARG GID=1000

RUN addgroup -S -g $GID $GROUP \
    && adduser -S -u $UID -D -G $GROUP $USER

COPY --from=quote-server-builder /ccnp-server/target/release/ccnp_server /bin
COPY --from=quote-server-builder /usr/bin/grpc-health-probe /usr/bin

USER $UID
CMD ["/bin/ccnp_server"]

