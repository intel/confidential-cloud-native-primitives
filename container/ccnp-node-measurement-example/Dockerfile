FROM python:3.12.1-alpine3.19

ARG USER=ccnp
ARG UID=1000
ARG GID=1000
ARG GROUP=ccnp

WORKDIR /run/ccnp-node-measurement-tool
RUN addgroup -S -g $GID $GROUP && adduser -S -u $UID -D -G $GROUP $USER
COPY container/ccnp-node-measurement-example/fetch_node_measurement.py ./

RUN apk update && apk add bash vim && pip install --no-cache-dir ccnp

USER $UID

ENTRYPOINT ["tail", "-f", "/dev/null"]
