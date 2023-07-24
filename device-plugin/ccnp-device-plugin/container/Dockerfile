From golang:1.20-alpine3.17 AS builder

RUN apk update \
    && apk add --no-cache protoc make

WORKDIR /usr/local/go/src/github.com/ccnp-device-plugin/
COPY . ./
RUN make


From alpine:3.17.0

WORKDIR /bin
COPY --from=builder /usr/local/go/src/github.com/ccnp-device-plugin/build/ccnp-device-plugin ./

CMD ["/bin/ccnp-device-plugin"]

