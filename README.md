# cenk1cenk2/beamer

[![pipeline status](https://gitlab.kilic.dev/docker/beamer/badges/master/pipeline.svg)](https://gitlab.kilic.dev/docker/beamer/-/commits/master) [![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/beamer)](https://hub.docker.com/repository/docker/cenk1cenk2/beamer) [![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/beamer)](https://hub.docker.com/repository/docker/cenk1cenk2/beamer) [![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/beamer)](https://hub.docker.com/repository/docker/cenk1cenk2/beamer) [![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/beamer)](https://github.com/cenk1cenk2/beamer)

## Description

Fast and dirty docker image to configure containers.

---

- [CLI Documentation](./CLI.md)

<!-- toc -->

<!-- tocstop -->

---

<!-- clidocs -->

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$VIZIER_CONFIG_FILE` | Configuration file to read from. | `String`<br/>`json(https://raw.githubusercontent.com/cenk1cenk2/docker-vizier/main/schema.json)` | `false` |  |

<!-- clidocsstop -->
