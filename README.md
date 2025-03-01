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
| `$BEAMER_MODE` | Mode to use. | `String`<br/>`enum(git)` | `false` | git |
| `$BEAMER_WORKING_DIRECTORY` | Working directory for cloning the data. | `String` | `false` | /tmp/beamer |

**Git**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$BEAMER_GIT_REPOSITORY` | Git repository to clone. | `String` | `false` |  |
| `$BEAMER_GIT_BRANCH` | Git branch to clone. | `String` | `false` |  |
| `$BEAMER_GIT_AUTH_METHOD` | Authentication method to use. | `String`<br/>`enum(none, ssh)` | `false` | none |
| `$BEAMER_GIT_PRIVATE_KEY` | Private key to use for SSH authentication. | `String` | `false` |  |
| `$BEAMER_GIT_PRIVATE_KEY_PASSWORD` | Password for the private key. | `String` | `false` |  |

<!-- clidocsstop -->
