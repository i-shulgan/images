FROM golang:1.27-bookworm AS go

COPY choose /choose

COPY rproxy /rproxy

RUN \
    cd /choose && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    cd /rproxy && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

FROM ubuntu:18.04

COPY --from=go /choose/choose /usr/local/bin/
COPY --from=go /rproxy/rproxy /usr/local/bin/

RUN \
    apt update && \
    apt -y install xvfb x11-utils x11vnc tigervnc-common tigervnc-viewer

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]

