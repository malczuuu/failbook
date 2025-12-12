# Releasing

[![Build & Push Docker Images](https://github.com/malczuuu/failbook/actions/workflows/docker-build.yml/badge.svg)](https://github.com/malczuuu/failbook/actions/workflows/docker-build.yml)
[![DockerHub](https://img.shields.io/docker/v/malczuuu/failbook?logo=docker&label=DockerHub)](https://hub.docker.com/r/malczuuu/failbook)

The project follows semantic versioning. A release is created by pushing an annotated git tag named `v1.2.3` with the message "Release 1.2.3". Preferably, use the `./tools/tagrelease` script, which ensures the tag is correctly formatted and prevents mistakes. Proper tag format is required to trigger build automation.

## Usage

See `./tools/tagrelease --help` for reference.

```txt
Create a git tag with a version number.

Arguments:  
• version — Version number in the format `1.2.3` or `1.2.3-suffix`

Options:  
• -h, --help — Show help message  
• --no-tag-prefix — Create a tag without the `v` prefix (`1.2.3` instead of `v1.2.3`)

Examples:  
• ./tools/tagrelease 1.2.3  
• ./tools/tagrelease 2.0.0-beta  
• ./tools/tagrelease --no-tag-prefix 1.2.4
```

## Example Release

To release version `1.2.3`, run:

```bash
./tools/tagrelease 1.2.3
```

On success, the script prints:

```txt
info: successfully created annotated tag 'v1.2.3', with message 'Release 1.2.3'
```

**Note:** You still need to push the created tag manually.

## Refreshing Images

For supported versions, images are automatically rebuild to update dependencies from base image.
Versions for rebuilding must be set in [`supported_versions`](./.github/utils/supported_versions).
