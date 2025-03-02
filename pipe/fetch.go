package pipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

			t.Pipe.Ctx.Fetch.State, err = json.Marshal(&GitStateFile{
				LastCommit: ref.Hash().String(),
			})
			if err != nil {
				return err
			}

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

			originalRef, err := r.Head()
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

			ref, err := r.Head()
			if err != nil {
				return err
			}

			t.Pipe.Ctx.Fetch.Dirty = true
			t.Pipe.Ctx.Fetch.State, err = json.Marshal(&GitStateFile{
				LastCommit: ref.Hash().String(),
			})
			if err != nil {
				return err
			}

			t.Log.Infof("Repository pulled successfully: %s -> %s", originalRef, ref)

			return nil
		})
}

func GitSyncDeletes(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("git", "sync", "delete").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Config.SyncDelete
			// || !t.Pipe.Ctx.Fetch.Dirty
		}).
		Set(func(t *Task[Pipe]) error {
			f, err := os.ReadFile(filepath.Join(t.Pipe.Config.TargetDirectory, t.Pipe.Config.StateFile))
			if errors.Is(err, os.ErrNotExist) {
				t.Log.Warnf("State file not found: %s", t.Pipe.Config.StateFile)

				return nil
			} else if err != nil {
				return err
			}
			state := &GitStateFile{}
			err = json.Unmarshal(f, &state)
			if err != nil {
				return err
			}

			t.Log.Infof("Syncing deleted files between commits: from %s", state.LastCommit)

			r := t.Pipe.Ctx.Git.Repository

			last, err := r.CommitObject(plumbing.NewHash(state.LastCommit))
			if err != nil {
				return err
			}
			head, err := r.Head()
			if err != nil {
				return err
			}
			now, err := r.CommitObject(head.Hash())
			if err != nil {
				return err
			}

			diff, err := last.Patch(now)
			if err != nil {
				return err
			}

			t.Log.Debugf("Diff: %v", diff.Stats().String())

			for _, s := range diff.Stats() {
				t.Log.Debugf("File: %v", s.Name)
			}

			return nil
		})
}
