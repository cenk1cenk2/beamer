package pipe

import (
	"time"

	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Pipe struct {
		Ctx    Ctx
		Config Config
		Git    Git
	}

	Config struct {
		Mode                       string `validate:"required,oneof=git"`
		WorkingDirectory           string
		TargetDirectory            string `validate:"required"`
		RootDirectory              string
		SyncDelete                 bool
		SyncDeleteEmptyDirectories bool
		StateFile                  string
		PullInterval               time.Duration
		IgnoreFile                 string
		ForceWorkflow              bool
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
				tl.JobSequence(
					tl.JobIf(
						func(_ floc.Context) bool {
							return tl.Pipe.Config.Mode == "git"
						},
						tl.JobSequence(
							GitConfigureAuthentication(tl).Job(),
							GitCloneRepository(tl).Job(),
							tl.JobLoopWithWaitAfter(
								tl.JobSequence(
									GitPull(tl).Job(),
									GitSyncDeletes(tl).Job(),
									Workflow(tl).Job(),
								),
								tl.Pipe.Config.PullInterval,
							),
						),
					),
				),
			)
		})
}
