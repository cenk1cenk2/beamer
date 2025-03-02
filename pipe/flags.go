package pipe

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/docker/beamer/internal/adapter"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONFIG = "Config"
)

var Flags = CombineFlags(
	[]cli.Flag{
		// category config

		&cli.StringFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "adapter",
			Usage:       fmt.Sprintf("Mode to use. enum(%v)", []string{ADAPTER_GIT}),
			Required:    false,
			Value:       ADAPTER_GIT,
			EnvVars:     []string{"BEAMER_ADAPTER"},
			Destination: &TL.Pipe.Config.Adapter,
		},

		&cli.DurationFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "interval",
			Usage:       "Interval between sync operations.",
			Required:    false,
			Value:       1 * time.Hour,
			EnvVars:     []string{"BEAMER_INTERVAL"},
			Destination: &TL.Pipe.Config.Interval,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "once",
			Usage:       "Run the workflow only once.",
			Required:    false,
			Value:       false,
			EnvVars:     []string{"BEAMER_ONCE"},
			Destination: &TL.Pipe.Config.Once,
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
			Name:        "working-directory",
			Usage:       "Working directory for cloning the data.",
			Required:    false,
			Value:       "/tmp/beamer",
			EnvVars:     []string{"BEAMER_WORKING_DIRECTORY"},
			Destination: &TL.Pipe.WorkingDirectory,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "root-directory",
			Usage:       "Root directory for the project.",
			Required:    false,
			Value:       "/",
			EnvVars:     []string{"BEAMER_ROOT_DIRECTORY"},
			Destination: &TL.Pipe.RootDirectory,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "target-directory",
			Usage:       "Target directory for the project.",
			Required:    true,
			Value:       "",
			EnvVars:     []string{"BEAMER_TARGET_DIRECTORY"},
			Destination: &TL.Pipe.TargetDirectory,
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
			Name:        "force-sync",
			Usage:       "Always force to sync the data, eventhough the state is not dirty.",
			Required:    false,
			Value:       false,
			EnvVars:     []string{"BEAMER_FORCE_SYNC"},
			Destination: &TL.Pipe.ForceSync,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "sync-delete",
			Usage:       "Delete files that are not in the source.",
			Required:    false,
			Value:       false,
			EnvVars:     []string{"BEAMER_SYNC_DELETE"},
			Destination: &TL.Pipe.SyncDelete,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "sync-delete-empty-directories",
			Usage:       "Delete empty directories after sync delete.",
			Required:    false,
			Value:       true,
			EnvVars:     []string{"BEAMER_SYNC_DELETE_EMPTY_DIRECTORIES"},
			Destination: &TL.Pipe.SyncDeleteEmptyDirectories,
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

		&cli.StringFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "file-comparator",
			Usage:       fmt.Sprintf("File comparator to use. enum(%v)", []string{comparator.COMPARATOR_SHA256, comparator.COMPARATOR_MD5}),
			Required:    false,
			Value:       comparator.COMPARATOR_MD5,
			EnvVars:     []string{"BEAMER_FILE_COMPARATOR"},
			Destination: &TL.Pipe.Config.FileComparator,
		},

		&cli.StringSliceFlag{
			Category: CATEGORY_CONFIG,
			Name:     "template-files",
			Usage:    "Template file extensions that should be rendered.",
			Required: false,
			Value:    cli.NewStringSlice(".tmpl", ".gotmpl"),
			EnvVars:  []string{"BEAMER_TEMPLATE_FILES"},
		},
	},
	adapter.GitAdapterFlags,
)

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.Config.TemplateFiles = tl.CliContext.StringSlice("template-files")

	return nil
}
