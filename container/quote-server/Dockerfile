FROM rust:1.71-alpine3.17 AS quote-server-builder

RUN apk update \
    && apk add --no-cache make protobuf-dev musl-dev wget

COPY service/quote-server /quote-server
COPY api /quote-server/api
RUN cd /quote-server && make build

# add rediness and liveness probe command
WORKDIR /usr/bin
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.18 && \
    wget -qO grpc-health-probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x grpc-health-probe

# ======================================================================================================================

From alpine:3.17.0

ARG USER=ccnp
ARG GROUP=ccnp
ARG UID=1000
ARG GID=1000

RUN addgroup -S -g $GID $GROUP \
    && adduser -S -u $UID -D -G $GROUP $USER

COPY --from=quote-server-builder /quote-server/target/release/quote_server /bin
COPY --from=quote-server-builder /usr/bin/grpc-health-probe /usr/bin

USER $USERNAME
CMD ["/bin/quote_server"]
