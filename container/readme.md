# [gesquive/paperless-uploader](https://github.com/gesquive/paperless-uploader)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/paperless-uploader/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/gesquive/paperless-uploader)
[![Build Status](https://img.shields.io/circleci/build/github/gesquive/paperless-uploader?style=flat-square)](https://circleci.com/gh/gesquive/paperless-uploader)
[![Coverage Report](https://img.shields.io/codecov/c/gh/gesquive/paperless-uploader?style=flat-square)](https://codecov.io/gh/gesquive/paperless-uploader)
[![Github Release](https://img.shields.io/github/v/tag/gesquive/paperless-uploader?style=flat-square)](https://github.com/gesquive/paperless-uploader)

# Supported Architectures

This image supports multiple architectures:

- `amd64`, `x86-64`
- `armv7`, `armhf`
- `arm64`, `aarch64`

Docker images are uploaded with using Docker manifest lists to make multi-platform deployments easer. More info can be found from [Docker](https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md#manifest-list)

You can simply pull the image using `gesquive/paperless-uploader` and docker should retreive the correct image for your architecture.

# Supported Tags
If you want a specific version of `paperless-uploader` you can pull it by specifying a version tag.

## Version Tags
This image provides versions that are available via tags. 

| Tag    | Description |
| ------ | ----------- |
| `latest` | Latest stable release |
| `0.9.0`  | Stable release v0.9.0 |
| `0.9.0-<git_hash>` | Development preview of version v0.9.0 at the given git hash |

# Usage

Here are some example snippets to help you get started creating a docker container.

## Docker CLI
```shell
docker run \
  --name=paperless-uploader \
  -v path/to/config:/config \
  -v path/to/watch:/watch \
  --restart unless-stopped \
  gesquive/paperless-uploader
```

## Docker Compose
Compatible with docker-compose v2 schemas.

```docker
---
version: "2"
services:
  paperless-uploader:
    image: gesquive/paperless-uploader
    container_name: paperless-uploader
    volumes:
      - path/to/config:/config
      - path/to/watch:/watch
    restart: unless-stopped
```

# Parameters
The container defines the following parameters that you can set:

| Parameter | Function |
| --------- | -------- |
| `-v /config`  | The paperless-uploader config goes here |
| `-v /watch`  | The directory to watch |
