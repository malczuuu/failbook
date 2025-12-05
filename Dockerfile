FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY . .

RUN go mod download

# builds a fully static Linux binary of Go application
# - CGO_ENABLED=0 disables CGO for full static linking
# - GOOS=linux targets Linux (cross-compilation if needed)
# - -a rebuilds all packages
# - -installsuffix cgo separates CGO-disabled build artifacts from any CGO-enabled ones
# - -ldflags '-extldflags "-static"' ensures the binary is statically linked (no libc dependency)
# - -o failbook sets the output binary name
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o failbook ./cmd/failbook

FROM alpine:3.19

RUN apk --no-cache add ca-certificates dumb-init

WORKDIR /app

COPY --from=builder /build/failbook .

COPY templates ./templates
COPY problem-docs ./problem-docs

# use non-root user
RUN addgroup -g 1000 failbook && \
    adduser -D -u 1000 -G failbook failbook && \
    chown -R failbook:failbook /app
USER failbook

EXPOSE 12001

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:12001/manage/health/live || exit 1

# use dumb-init to handle signals properly
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["./failbook"]
