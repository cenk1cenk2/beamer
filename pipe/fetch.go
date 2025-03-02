package pipe

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

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
			t.Log.Infof("Cloning repository: %s@%s", tl.Pipe.Git.Repository, tl.Pipe.Git.Branch)

			r, err := git.PlainClone(t.Pipe.Config.WorkingDirectory, false, &git.CloneOptions{
				URL:           tl.Pipe.Git.Repository,
				Progress:      t.Log.Writer(),
				Auth:          t.Pipe.Ctx.Git.AuthMethod,
				SingleBranch:  true,
				ReferenceName: plumbing.ReferenceName(t.Pipe.Git.Branch),
			})
			if errors.Is(err, git.ErrRepositoryAlreadyExists) {
				t.Log.Warnf("Repository already exists. Skipping clone.")

				r, err = git.PlainOpen(t.Pipe.Config.WorkingDirectory)
				if err != nil {
					return err
				}

				remote, err := r.Remote("origin")
				if err != nil {
					return err
				}

				if remote.Config().URLs[0] != tl.Pipe.Git.Repository {
					return fmt.Errorf("Remote repository URL does not match the provided URL: %s -> %s", remote.Config().URLs[0], tl.Pipe.Git.Repository)
				}
			} else if err != nil {
				return err
			}

			t.Pipe.Ctx.Git.Repository = r

			ref, err := r.Head()
			if err != nil {
				return err
			}
			t.Log.Infof("Repository cloned successfully: %s", ref)

			return nil
		})
}

func GitPull(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("git", "pull").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof("Pulling repository: %s@%s", tl.Pipe.Git.Repository, tl.Pipe.Git.Branch)

			r := t.Pipe.Ctx.Git.Repository
			w, err := r.Worktree()
			if err != nil {
				return err
			}

			err = w.Pull(&git.PullOptions{
				Progress:      t.Log.Writer(),
				Auth:          t.Pipe.Ctx.Git.AuthMethod,
				ReferenceName: plumbing.ReferenceName("HEAD"),
			})
			if errors.Is(err, git.NoErrAlreadyUpToDate) {
				ref, err := r.Head()
				if err != nil {
					return err
				}

				t.Log.Infof("Repository is already up-to-date: %s", ref)

				return nil
			} else if err != nil {
				return err
			}

			t.Pipe.Ctx.Fetch.Dirty = true

			return nil
		})
}
