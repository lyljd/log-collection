FROM alpine

LABEL maintainer="ljd9031@163.com"
LABEL version="1.0"

RUN apk update && apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

ENV GIN_MODE="release"

WORKDIR /app

COPY log-configuration .
COPY ui.html .
COPY static static

ENTRYPOINT ./log-configuration
