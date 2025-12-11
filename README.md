# Failbook

[![Go Build](https://github.com/malczuuu/failbook/actions/workflows/go-build.yml/badge.svg)](https://github.com/malczuuu/failbook/actions/workflows/go-build.yml)
[![DockerHub](https://img.shields.io/docker/v/malczuuu/failbook?label=DockerHub)](https://hub.docker.com/r/malczuuu/failbook)
[![License](https://img.shields.io/github/license/malczuuu/failbook)](https://github.com/malczuuu/failbook/blob/main/LICENSE)

A simple HTTP API error documentation service written in Go. Failbook serves markdown-powered error documentation pages
for HTTP API error responses.

While working on [Problem4J](https://github.com/malczuuu/problem4j-spring), the initial assumption was that the `type`
field must be a resolvable HTTP URI. It turns out this is not always the case, but for the purpose of experimentation
this application was created. It allows configuring static `Problem` documentation pages.

## Quick Start

### Using Go

Current instructions use the [Taskfile](https://taskfile.dev/docs/getting-started) tool.

```bash
git clone https://github.com/malczuuu/failbook.git
cd failbook

task test
task build

./dist/failbook
```

Visit http://localhost:12001 to see your error documentation.

### Using Docker

Docker images are available on Docker Hub as [`malczuuu/failbook`](https://hub.docker.com/r/malczuuu/failbook).

```bash
docker run -p 12001:12001 -v $(pwd)/problem-docs:/failbook/problem-docs:ro malczuuu/failbook:latest
```

See also `failbook-compose/` for a `docker-compose.yaml` demo.

To build the Docker image on a local machine, use the `build-docker` task. It will also build the
`failbook-builder:latest` image to speed up subsequent builds. See `Dockerfile.builder` for more details.

```bash
task build-docker
```

How to do this without `task`:

<details>
<summary><b>Expand...</b></summary>

```bash
docker build -f Dockerfile.builder -t failbook-builder:latest .
docker build -f Dockerfile -t malczuuu/failbook:latest .
```

</details>

## Configuration

Failbook is configured via environment variables:

| Variable                      | Default                  | Description                                                    |
|-------------------------------|--------------------------|----------------------------------------------------------------|
| `FAILBOOK_PORT`               | `12001`                  | HTTP server port                                               |
| `FAILBOOK_LOG_LEVEL`          | `info`                   | Log level (`trace`, `debug`, `info`, `warn`, `error`, `fatal`) |
| `FAILBOOK_HEALTH_ENABLED`     | `false`                  | Enable health check endpoints                                  |
| `FAILBOOK_PROMETHEUS_ENABLED` | `false`                  | Enable Prometheus metrics endpoint                             |
| `FAILBOOK_PROBLEM_DOCS_DIR`   | `/failbook/problem-docs` | Directory containing error YAML files                          |
| `FAILBOOK_BASE_HREF`          | (empty)                  | Base path for reverse proxy deployments (e.g., `/api/docs`)    |

### Example

```bash
export FAILBOOK_PORT=8080
export FAILBOOK_LOG_LEVEL=debug
export FAILBOOK_PROMETHEUS_ENABLED=true
export FAILBOOK_BASE_HREF=/api/docs
./failbook
```

## Error Configuration Format

Error documentation is defined in YAML files in the `errors/` directory. Each file may contain one or more error
definitions.

### YAML Schema

```yaml
version: "1"               # Required: Schema version, must be "1"
id: "404"                  # Required: Unique error identifier
name: "Validation Failed"  # Optional: Composed as "{title} {status_code}" if not provided
title: "Not Found"         # Required: Short error title
status_code: 404           # Required: HTTP status code
summary: "The requested resource could not be found"  # Required: Brief summary (shown on index)
description: |             # Required: Detailed description (supports Markdown)
  ## What Happened
  
  The server cannot find the requested resource.
  
  ## Common Causes
  
  - **Invalid URL**: The URL may be misspelled
  - **Deleted Resource**: The resource may have been removed
  - **Moved Resource**: The resource may have been moved
  
  ## What To Do
  
  1. Check the URL for typos
  2. Verify the resource exists
  3. Check the API documentation

links:                 # Optional: Related links
  - title: "API Documentation"
    url: "https://api.example.com/docs"
  - title: "Support"
    url: "https://support.example.com"
```

### Multi-Document YAML Files

You can define multiple errors in a single YAML file using document separators (`---`):

```yaml
version: "1"
id: "400"
title: "Bad Request"
status_code: 400
summary: "The request was malformed"
description: "Your request contains invalid syntax."
---
version: "1"
id: "401"
title: "Unauthorized"
status_code: 401
summary: "Authentication is required"
description: "You must authenticate to access this resource."
---
version: "1"
id: "403"
title: "Forbidden"
status_code: 403
summary: "Access denied"
description: "You don't have permission to access this resource."
```

### Markdown Support

The `description` field supports Markdown, powered by the [`yuin/goldmark`](https://github.com/yuin/goldmark) library.

## Endpoints

### Application Endpoints

- `GET /` — error documentation index page  
- `GET /:id` — individual error detail page (`id` may contain multiple path segments)

### Management Endpoints

- `GET /manage/health/live` — liveness probe (always returns 200 OK, if enabled)  
- `GET /manage/health/ready` — readiness probe (returns 200 when ready, 503 when not, if enabled)  
- `GET /manage/prometheus` — Prometheus metrics (if enabled)
