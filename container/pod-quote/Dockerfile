FROM rust:1.74.0 as pod-quote-builder

RUN apt-get update && apt-get install -y \
    build-essential \
    libprotobuf-dev \
    protobuf-compiler \
    libssl-dev \
    pkg-config

WORKDIR /app

COPY . .

COPY service/pod-quote /pod-quote

RUN cd /pod-quote && make build

FROM rust:1.74.0-slim

RUN apt-get update && apt-get install -y curl

WORKDIR /app
COPY --from=pod-quote-builder /pod-quote/target/release/pod_quote /app/pod_quote

RUN chmod a+x /app/pod_quote

# Run the sleep command for demonstration purposes
CMD ["/app/pod_quote"]
