# Failbook

[![Go Build with Taskfile](https://github.com/malczuuu/failbook/actions/workflows/go-build-with-taskfile.yml/badge.svg)](https://github.com/malczuuu/failbook/actions/workflows/go-build-with-taskfile.yml)
[![License](https://img.shields.io/github/license/malczuuu/failbook)](https://github.com/malczuuu/failbook/blob/main/LICENSE)

⚠️ **Project Status: DRAFT / WORK IN PROGRESS** ⚠️

A simple HTTP API error documentation service written in Go. Failbook serves simple, markdown-powered error 
documentation pages for HTTP API error responses.

While writing [Problem4J](https://github.com/malczuuu/problem4j-spring) library, my assumption was that `type` field
must be a resolvable HTTP URI. This turns out is not always the case, but for the sole purpose of experimentation,
this application was created. It allows configuring a simple static `Problem` documentation pages.

## Quick Start

### Using Go (via Taskfile)

```bash
git clone https://github.com/malczuuu/failbook.git
cd failbook

task test
task build

./dist/failbook
```

Visit http://localhost:12001 to see your error documentation.

### Using Docker

```bash
docker build -t failbook .
docker run -p 12001:12001 -v $(pwd)/problem-docs:/app/problem-docs:ro failbook
```

## Configuration

Failbook is configured via environment variables:

| Variable                      | Default          | Description                                                 |
|-------------------------------|------------------|-------------------------------------------------------------|
| `FAILBOOK_PORT`               | `12001`          | HTTP server port                                            |
| `FAILBOOK_LOG_LEVEL`          | `info`           | log level (trace, debug, info, warn, error, fatal)          |
| `FAILBOOK_PROMETHEUS_ENABLED` | `false`          | enable Prometheus metrics endpoint                          |
| `FAILBOOK_PROBLEM_DOCS_DIR`   | `./problem-docs` | directory containing error YAML files                       |
| `FAILBOOK_BASE_HREF`          | (empty)          | base path for reverse proxy deployments (e.g., `/api/docs`) |

### Example

```bash
export FAILBOOK_PORT=8080
export FAILBOOK_LOG_LEVEL=debug
export FAILBOOK_PROMETHEUS_ENABLED=true
export FAILBOOK_BASE_HREF=/api/docs
./failbook
```

## Error Configuration Format

Error documentation is defined in YAML files in the `errors/` directory. Each file can contain one or more error
definitions.

### YAML Schema

```yaml
version: "1"           # Required: Schema version, must be "1"
id: "404"              # Required: Unique error identifier
title: "Not Found"     # Required: Short error title
status_code: 404       # Required: HTTP status code
summary: "The requested resource could not be found"  # Required: Brief summary (shown on index)
description: |         # Required: Detailed description (supports Markdown)
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

The `description` field supports Markdown, powered by [`yuin/goldmark`](https://github.com/yuin/goldmark) library.

## Endpoints

### Application Endpoints

- `GET /` - error documentation index page
- `GET /:id` - individual error detail page

### Management Endpoints

- `GET /manage/health/live` - liveness probe (always returns 200 OK)
- `GET /manage/health/ready` - readiness probe (returns 200 when ready, 503 when not)
- `GET /manage/prometheus` - prometheus metrics (when enabled)
