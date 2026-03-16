FROM golang:1.25-alpine

WORKDIR /build

RUN apk add --no-cache curl git
RUN curl -sL https://taskfile.dev/install.sh | sh
