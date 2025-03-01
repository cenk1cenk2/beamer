package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Pipe struct {
		Ctx    Ctx
		Config Config
		Git    Git
	}

	Config struct {
		Mode             string `validate:"required,oneof=git"`
		WorkingDirectory string
	}

	Git struct {
		AuthMethod            string `validate:"oneof=none ssh"`
		Repository            string `validate:"required"`
		Branch                string
		SshPrivateKey         string `validate:"required_if=Inner.AuthMethod ssh"`
		SshPrivateKeyPassword string
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(1).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				WorkflowGit(tl).Job(),
			)
		})
}
