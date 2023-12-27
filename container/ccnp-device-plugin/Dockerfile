From golang:1.20-alpine3.18 AS builder

RUN apk update \
    && apk add --no-cache protoc make

WORKDIR /usr/local/go/src/github.com/ccnp-device-plugin/
COPY device-plugin/ccnp-device-plugin ./
RUN make


From alpine:3.18.5

WORKDIR /bin
COPY --from=builder /usr/local/go/src/github.com/ccnp-device-plugin/build/ccnp-device-plugin ./

CMD ["/bin/ccnp-device-plugin"]

