package pipe

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func WorkflowGit(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("git").
		SetJobWrapper(func(_ Job, _ *Task[Pipe]) Job {
			return tl.JobSequence(
				GitConfigureAuthentication(tl).Job(),
				GitCloneRepository(tl).Job(),
			)
		})
}

func GitConfigureAuthentication(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("git", "conf", "auth").
		Set(func(t *Task[Pipe]) error {
			switch t.Pipe.Git.AuthMethod {
			case "ssh":
				t.Log.Infof("Using SSH authentication for git.")

				am, err := ssh.NewPublicKeys("git", t.Pipe.Ctx.Git.SshPrivateKey, t.Pipe.Git.SshPrivateKeyPassword)
				if err != nil {
					return err
				}

				t.Pipe.Ctx.Git.AuthMethod = am
			}

			return nil
		})
}

func GitCloneRepository(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("git", "clone").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof("Cloning repository: %s", tl.Pipe.Git.Repository)

			repository, err := git.PlainClone(tl.Pipe.Config.WorkingDirectory, true, &git.CloneOptions{
				URL:           t.Pipe.Git.Repository,
				SingleBranch:  true,
				Depth:         1,
				ReferenceName: plumbing.ReferenceName(t.Pipe.Git.Branch),
				Progress:      t.Log.Writer(),
				Auth:          t.Pipe.Ctx.Git.AuthMethod,
			})
			if err != nil {
				return err
			}

			ref, err := repository.Head()
			if err != nil {
				return err
			}
			t.Log.Infof("Repository cloned successfully: %s", ref)

			return nil
		})
}
