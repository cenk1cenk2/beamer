package pipe

import (
	"encoding/base64"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONFIG = "Config"
	CATEGORY_GIT    = "Git"
)

var Flags = []cli.Flag{
	// category config

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "mode",
		Usage:       "Mode to use. enum(git)",
		Required:    false,
		Value:       "git",
		EnvVars:     []string{"BEAMER_MODE"},
		Destination: &TL.Pipe.Config.Mode,
	},

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "working-directory",
		Usage:       "Working directory for cloning the data.",
		Required:    false,
		Value:       "/tmp/beamer",
		EnvVars:     []string{"BEAMER_WORKING_DIRECTORY"},
		Destination: &TL.Pipe.Config.WorkingDirectory,
	},

	&cli.DurationFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "pull-interval",
		Usage:       "Interval to wait between pull operations.",
		Required:    false,
		Value:       time.Second * 5,
		EnvVars:     []string{"BEAMER_PULL_INTERVAL"},
		Destination: &TL.Pipe.Config.PullInterval,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "force-workflow",
		Usage:       "Force workflow to run even if the data is not dirty.",
		Required:    false,
		Value:       false,
		EnvVars:     []string{"BEAMER_FORCE_WORKFLOW"},
		Destination: &TL.Pipe.Config.ForceWorkflow,
	},

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "root-directory",
		Usage:       "Root directory for the project.",
		Required:    false,
		Value:       "/",
		EnvVars:     []string{"BEAMER_ROOT_DIRECTORY"},
		Destination: &TL.Pipe.Config.RootDirectory,
	},

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "target-directory",
		Usage:       "Target directory for the project.",
		Required:    true,
		Value:       "",
		EnvVars:     []string{"BEAMER_TARGET_DIRECTORY"},
		Destination: &TL.Pipe.Config.TargetDirectory,
	},

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "ignore-file",
		Usage:       "File to use for ignoring files.",
		Required:    false,
		Value:       ".beamer-ignore",
		EnvVars:     []string{"BEAMER_IGNORE_FILE"},
		Destination: &TL.Pipe.Config.IgnoreFile,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "sync-delete",
		Usage:       "Delete files that are not in the source.",
		Required:    false,
		Value:       false,
		EnvVars:     []string{"BEAMER_SYNC_DELETE"},
		Destination: &TL.Pipe.Config.SyncDelete,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "sync-delete-empty-directories",
		Usage:       "Delete empty directories after sync delete.",
		Required:    false,
		Value:       true,
		EnvVars:     []string{"BEAMER_SYNC_DELETE_EMPTY_DIRECTORIES"},
		Destination: &TL.Pipe.Config.SyncDeleteEmptyDirectories,
	},

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "state-file",
		Usage:       "File to use for storing state.",
		Required:    false,
		Value:       ".beamer",
		EnvVars:     []string{"BEAMER_STATE_FILE"},
		Destination: &TL.Pipe.Config.StateFile,
	},

	// category git

	&cli.StringFlag{
		Category:    CATEGORY_GIT,
		Name:        "git-repository",
		Usage:       "Git repository to clone.",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"BEAMER_GIT_REPOSITORY"},
		Destination: &TL.Pipe.Git.Repository,
	},

	&cli.StringFlag{
		Category:    CATEGORY_GIT,
		Name:        "git-branch",
		Usage:       "Git branch to clone.",
		Required:    false,
		Value:       "HEAD",
		EnvVars:     []string{"BEAMER_GIT_BRANCH"},
		Destination: &TL.Pipe.Git.Branch,
	},

	&cli.StringFlag{
		Category:    CATEGORY_GIT,
		Name:        "git-auth-method",
		Usage:       "Authentication method to use. enum(none, ssh)",
		Required:    false,
		Value:       "none",
		EnvVars:     []string{"BEAMER_GIT_AUTH_METHOD"},
		Destination: &TL.Pipe.Git.AuthMethod,
	},

	&cli.StringFlag{
		Category:    CATEGORY_GIT,
		Name:        "git-private-key",
		Usage:       "Private key to use for SSH authentication.",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"BEAMER_GIT_SSH_PRIVATE_KEY"},
		Destination: &TL.Pipe.Git.SshPrivateKey,
	},

	&cli.StringFlag{
		Category:    CATEGORY_GIT,
		Name:        "git-private-key-password",
		Usage:       "Password for the private key.",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"BEAMER_GIT_PRIVATE_KEY_PASSWORD"},
		Destination: &TL.Pipe.Git.SshPrivateKeyPassword,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	switch tl.Pipe.Git.AuthMethod {
	case "ssh":
		k := tl.CliContext.String("git-private-key")
		if stat, err := os.Stat(k); err == nil && !stat.IsDir() {
			key, err := os.ReadFile(k)
			if err != nil {
				return err
			}

			tl.Log.Debug("Using SSH private key as file.")

			tl.Pipe.Ctx.Git.SshPrivateKey = key
		}

		key, err := base64.StdEncoding.DecodeString(k)
		if err != nil {
			return err
		}

		tl.Log.Debug("Using SSH private key directly from flag.")

		tl.Pipe.Ctx.Git.SshPrivateKey = key
	}

	return nil
}
