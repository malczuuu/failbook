FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY . .

RUN chmod +x ./docker/healthcheck.sh
RUN apk add --no-cache curl git
RUN curl -sL https://taskfile.dev/install.sh | sh

RUN ./bin/task build-prod

FROM alpine:3.19

ENV FAILBOOK_PORT=12001
ENV FAILBOOK_LOG_LEVEL=info
ENV FAILBOOK_HEALTH_ENABLED=false
ENV FAILBOOK_PROMETHEUS_ENABLED=false
ENV FAILBOOK_PROBLEM_DOCS_DIR=./problem-docs
ENV FAILBOOK_BASE_HREF=

RUN apk --no-cache add ca-certificates dumb-init

WORKDIR /failbook

COPY --from=builder /build/dist/failbook .
COPY --from=builder /build/docker/healthcheck.sh .

COPY templates ./templates
COPY problem-docs ./problem-docs

# use non-root user
RUN addgroup -g 1000 failbook && \
    adduser -D -u 1000 -G failbook failbook && \
    chown -R failbook:failbook /failbook
USER failbook

EXPOSE 12001

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD /failbook/healthcheck.sh || exit 1

# use dumb-init to handle signals properly
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/failbook/failbook"]
