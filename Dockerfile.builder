FROM golang:1.24-alpine

WORKDIR /build

RUN apk add --no-cache curl git
RUN curl -sL https://taskfile.dev/install.sh | sh
