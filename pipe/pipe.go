package pipe

import (
	"time"

	"github.com/workanator/go-floc/v3"
	"gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/adapter"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Pipe struct {
		internal.ServiceFlags

		Ctx    Ctx
		Config Config
	}

	Config struct {
		Adapter        Adapter `validate:"required,oneof=git"`
		Once           bool
		StateFile      string
		LockFile       string
		Interval       time.Duration
		IgnoreFile     string
		ForceWorkflow  bool
		FileComparator comparator.Comparator `validate:"oneof=sha256 md5"`
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}
var a adapter.Adapter

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(2).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			if err := ProcessFlags(tl); err != nil {
				return err
			}

			return tl.RunJobs(Setup(tl).Job())
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			jobs := tl.JobIf(func(_ floc.Context) bool {
				return !tl.Pipe.Ctx.LockFile.IsLocked()
			},
				tl.JobSequence(
					tl.CreateBasicJob(func() error {
						return tl.Pipe.Ctx.LockFile.Lock()
					}),
					a.Sync(),
					Workflow(tl).Job(),
					a.Finalize(),
					tl.GuardAlways(
						tl.CreateBasicJob(func() error {
							return tl.Pipe.Ctx.LockFile.Unlock()
						}),
					),
				),
				tl.CreateBasicJob(func() error {
					tl.Log.Warnf("Another process is running. Skipping this run.")

					return nil
				}),
			)

			return tl.JobSequence(
				a.Init(),
				tl.JobIf(
					func(_ floc.Context) bool {
						return tl.Pipe.Config.Once
					},
					tl.JobThen(jobs),
					tl.JobElse(
						tl.JobLoopWithWaitAfter(
							tl.JobSequence(
								jobs,
							),
							tl.Pipe.Config.Interval,
						),
					),
				),
			)
		})
}
