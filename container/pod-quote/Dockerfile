FROM rust:1.74.0 as pod-quote-builder

ENV HTTP_PROXY=child-prc.intel.com:913
ENV HTTPS_PROXY=child-prc.intel.com:913
ENV FTP_PROXY=child-prc.intel.com:913

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    make \
    libprotobuf-dev \
    protobuf-compiler \
    musl-dev \
    wget \
    libssl-dev \
    pkg-config

WORKDIR /app

COPY . .

COPY service/pod-quote /pod-quote

RUN cd /pod-quote && make build

FROM rust:1.74.0

WORKDIR /app
COPY --from=pod-quote-builder /pod-quote/target/release/pod_quote /app/pod_quote

RUN chmod a+x /app/pod_quote

# Run the sleep command for demonstration purposes
CMD ["/app/pod_quote"]
