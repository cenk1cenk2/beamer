# beamer

Beam a configuration up to a container.

`beamer [FLAGS]`

## Flags

**Adapter Git**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$BEAMER_GIT_REPOSITORY` | Git repository to clone. | `String` | `false` |  |
| `$BEAMER_GIT_BRANCH` | Git branch to clone. | `String` | `false` | HEAD |
| `$BEAMER_GIT_AUTH_METHOD` | Authentication method to use. | `String`<br/>`enum(none, ssh)` | `false` | none |
| `$BEAMER_GIT_SSH_PRIVATE_KEY` | Private key to use for SSH authentication. | `String` | `false` |  |
| `$BEAMER_GIT_PRIVATE_KEY_PASSWORD` | Password for the private key. | `String` | `false` |  |

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `String`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |

**Config**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$BEAMER_ADAPTER` | Mode to use. | `String`<br/>`enum([git])` | `false` | git |
| `$BEAMER_INTERVAL` | Interval between sync operations. | `Duration` | `false` | 1h0m0s |
| `$BEAMER_ONCE` | Run the workflow only once. | `Bool` | `false` | false |
| `$BEAMER_FORCE_WORKFLOW` | Force workflow to run even if the data is not dirty. | `Bool` | `false` | false |
| `$BEAMER_WORKING_DIRECTORY` | Working directory for cloning the data. | `String` | `false` | /tmp/beamer |
| `$BEAMER_ROOT_DIRECTORY` | Root directory for the project. | `String` | `false` | / |
| `$BEAMER_TARGET_DIRECTORY` | Target directory for the project. | `String` | `true` |  |
| `$BEAMER_IGNORE_FILE` | File to use for ignoring files. | `String` | `false` | .beamer-ignore |
| `$BEAMER_FORCE_SYNC` | Always force to sync the data, eventhough the state is not dirty. | `Bool` | `false` | false |
| `$BEAMER_SYNC_DELETE` | Delete files that are not in the source. | `Bool` | `false` | false |
| `$BEAMER_SYNC_DELETE_EMPTY_DIRECTORIES` | Delete empty directories after sync delete. | `Bool` | `false` | true |
| `$BEAMER_STATE_FILE` | File to use for storing state. | `String` | `false` | .beamer |
| `$BEAMER_FILE_COMPARATOR` | File comparator to use. | `String`<br/>`enum([sha256 md5])` | `false` | md5 |
| `$BEAMER_TEMPLATE_FILES` | Template file extensions that should be rendered. | `StringSlice` | `false` | ".tmpl", ".gotmpl" |
