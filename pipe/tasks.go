package pipe

import (
	"fmt"
	"path/filepath"

	"gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/adapter"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
	"gitlab.kilic.dev/docker/beamer/internal/operations"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("setup").
		Set(func(t *Task[Pipe]) error {
			ctx := &internal.ServiceCtx{
				Log:              t.Log,
				WorkingDirectory: t.Pipe.WorkingDirectory,
				TargetDirectory:  t.Pipe.TargetDirectory,
				RootDirectory:    t.Pipe.RootDirectory,
				Flags: &internal.ServiceFlags{
					ForceSync:                  t.Pipe.ForceSync,
					SyncDelete:                 t.Pipe.SyncDelete,
					SyncDeleteEmptyDirectories: t.Pipe.SyncDeleteEmptyDirectories,
					TemplateFiles:              t.Pipe.TemplateFiles,
				},
			}
			ctx.State = internal.NewState(ctx, filepath.Join(t.Pipe.TargetDirectory, t.Pipe.Config.StateFile))
			t.Pipe.Ctx.State = ctx.State

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

				t.Log.Infof("Using SHA256 file comparator.")
			case comparator.COMPARATOR_MD5:
				t.Pipe.Ctx.FileComparator = comparator.NewFileComparatorMd5()

				t.Log.Infof("Using MD5 file comparator.")
			default:
				return fmt.Errorf("File comparator %s is not supported", t.Pipe.Config.FileComparator)
			}

			t.Pipe.Ctx.LockFile = operations.NewLockFile(
				t.Log.WithField(LOG_FIELD_CONTEXT, "locker"),
				t.Pipe.TargetDirectory,
				t.Pipe.Config.LockFile,
			)
			t.Log.Debugf("Lock file: %s", t.Pipe.Ctx.LockFile.Path())

			return nil
		})
}
