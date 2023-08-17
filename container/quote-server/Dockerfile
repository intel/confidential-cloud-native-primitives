FROM ubuntu:22.04 AS quote-server-builder

RUN apt-get update && apt-get install -y git golang curl build-essential  clang protobuf-compiler protobuf-c-compiler libprotobuf-c-dev libprotobuf-c1

RUN curl --proto '=https' --tlsv1.3 https://sh.rustup.rs -sSf > rustup.sh
RUN /bin/sh ./rustup.sh --profile minimal -y && rm ./rustup.sh

COPY service/quote-server /quote-server
COPY api /quote-server/api
RUN cd /quote-server && make build

# add rediness and liveness probe command
RUN git clone https://github.com/grpc-ecosystem/grpc-health-probe.git
RUN cd /grpc-health-probe && go build

# ======================================================================================================================

FROM ubuntu:22.04 as quote-server-base

ARG USERNAME=ccnp
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

COPY --from=quote-server-builder /quote-server/target/release/quote_server /bin
COPY --from=quote-server-builder /grpc-health-probe/grpc-health-probe /usr/bin

USER $USERNAME
CMD ["/bin/quote_server"]
