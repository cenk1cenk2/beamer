# beamer

Beam a configuration up to a container.

`beamer [FLAGS]`

## Flags

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
