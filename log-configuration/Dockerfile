FROM alpine

LABEL maintainer="ljd9031@163.com"
LABEL version="1.0"

ENV GIN_MODE="release"

WORKDIR /app

COPY log-configuration .
COPY ui.html .
COPY static static

ENTRYPOINT ./log-configuration
