package pipe

import (
	"fmt"
	"path/filepath"

	services "gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/adapter"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("setup").
		Set(func(t *Task[Pipe]) error {
			ctx := &services.ServiceCtx{
				Log:              t.Log,
				WorkingDirectory: t.Pipe.WorkingDirectory,
				TargetDirectory:  t.Pipe.TargetDirectory,
				RootDirectory:    t.Pipe.RootDirectory,
				Flags: &services.ServiceFlags{
					ForceSync:                  t.Pipe.ForceSync,
					SyncDelete:                 t.Pipe.SyncDelete,
					SyncDeleteEmptyDirectories: t.Pipe.SyncDeleteEmptyDirectories,
				},
			}
			ctx.State = services.NewState(ctx, filepath.Join(t.Pipe.TargetDirectory, t.Pipe.Config.StateFile))

			var err error

			switch tl.Pipe.Config.Adapter {
			case ADAPTER_GIT:
				a, err = adapter.NewGitAdapter(tl.Plumber, ctx)
				if err != nil {
					return err
				}

				t.Log.Infof("Using Git adapter.")
			default:
				return fmt.Errorf("Adapter %s is not supported", tl.Pipe.Config.Adapter)
			}

			switch t.Pipe.Config.FileComparator {
			case comparator.COMPARATOR_SHA256:
				t.Pipe.Ctx.FileComparator = comparator.NewFileComparatorSha256()

				t.Log.Debugf("Using SHA256 file comparator.")
			}

			return nil
		})
}
