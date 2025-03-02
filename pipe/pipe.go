package pipe

import (
	"time"

	services "gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/adapter"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Pipe struct {
		services.ServiceFlags

		Ctx    Ctx
		Config Config
	}

	Config struct {
		Adapter        Adapter `validate:"required,oneof=git"`
		StateFile      string
		Interval       time.Duration
		IgnoreFile     string
		ForceWorkflow  bool
		FileComparator comparator.Comparator `validate:"oneof=sha256"`
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
			return ProcessFlags(tl)
		}).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return tl.RunJobs(Setup(tl).Job())
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				a.Init(),
				tl.JobLoopWithWaitAfter(
					tl.JobSequence(
						a.Sync(),
						Workflow(tl).Job(),
						a.Finalize(),
					),
					tl.Pipe.Config.Interval,
				),
			)
		})
}
