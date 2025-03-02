package adapter

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/operations"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type BeamerGitAuthMethod = string

const (
	BEAMER_GIT_AUTH_METHOD_SSH BeamerGitAuthMethod = "ssh"
)

const category_git = "Adapter Git"

var gitAdapterFlags = &GitAdapterConfig{}

var GitAdapterFlags = []cli.Flag{
	&cli.StringFlag{
		Category:    category_git,
		Name:        "git-repository",
		Usage:       "Git repository to clone.",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"BEAMER_GIT_REPOSITORY"},
		Destination: &gitAdapterFlags.Repository,
	},

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git-branch",
		Usage:       "Git branch to clone.",
		Required:    false,
		Value:       "HEAD",
		EnvVars:     []string{"BEAMER_GIT_BRANCH"},
		Destination: &gitAdapterFlags.Branch,
	},

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git-auth-method",
		Usage:       "Authentication method to use. enum(none, ssh)",
		Required:    false,
		Value:       "none",
		EnvVars:     []string{"BEAMER_GIT_AUTH_METHOD"},
		Destination: &gitAdapterFlags.AuthMethod,
	},

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git-private-key",
		Usage:       "Private key to use for SSH authentication.",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"BEAMER_GIT_SSH_PRIVATE_KEY"},
		Destination: &gitAdapterFlags.SshPrivateKey,
	},

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git-private-key-password",
		Usage:       "Password for the private key.",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"BEAMER_GIT_PRIVATE_KEY_PASSWORD"},
		Destination: &gitAdapterFlags.SshPrivateKeyPassword,
	},
}

type GitAdapterConfig struct {
	Repository            string `validate:"required"`
	Branch                string
	AuthMethod            BeamerGitAuthMethod `validate:"oneof=none ssh"`
	SshPrivateKey         string              `validate:"required_if=Inner.AuthMethod ssh"`
	SshPrivateKeyPassword string
}

type GitAdapter struct {
	ctx              *internal.ServiceCtx
	tl               *TaskList[any]
	repository       *git.Repository
	authMethod       transport.AuthMethod
	config           *GitAdapterConfig
	workingDirectory string
	state            *GitAdapterState
}

type GitAdapterState struct {
	LastCommit string `json:"lastCommit"`
}

var _ Adapter = (*GitAdapter)(nil)

func NewGitAdapter(p *Plumber, ctx *internal.ServiceCtx) (*GitAdapter, error) {
	log := ctx.Log.WithField(LOG_FIELD_CONTEXT, "adapter").WithField(LOG_FIELD_STATUS, "git")

	adapter := &GitAdapter{
		ctx:              ctx,
		config:           gitAdapterFlags,
		workingDirectory: ctx.WorkingDirectory,
		tl:               (&TaskList[any]{}).New(p),
	}

	switch gitAdapterFlags.AuthMethod {
	case BEAMER_GIT_AUTH_METHOD_SSH:
		var key []byte

		if f := operations.NewFile(gitAdapterFlags.SshPrivateKey); f.IsFile() {
			k, err := f.ReadFile()
			if err != nil {
				return nil, err
			}
			key = k

			log.Debug("Using SSH private key as file.")
		} else {
			k, err := base64.StdEncoding.DecodeString(gitAdapterFlags.SshPrivateKey)
			if err != nil {
				return nil, err
			}
			key = k

			log.Debug("Using SSH private key directly from flag.")
		}

		log.Infof("Using SSH authentication for git.")

		am, err := ssh.NewPublicKeys("git", key, gitAdapterFlags.SshPrivateKeyPassword)
		if err != nil {
			return nil, err
		}
		adapter.authMethod = am
	}

	return adapter, nil
}

func (a *GitAdapter) Init() Job {
	return a.tl.CreateTask("init").
		Set(func(t *Task[any]) error {
			t.Log.Infof("Cloning repository: %s@%s", a.config.Repository, a.config.Branch)

			r, err := git.PlainClone(a.workingDirectory, false, &git.CloneOptions{
				URL:           a.config.Repository,
				Progress:      t.Log.Writer(),
				Auth:          a.authMethod,
				SingleBranch:  true,
				ReferenceName: plumbing.ReferenceName(a.config.Branch),
			})
			if errors.Is(err, git.ErrRepositoryAlreadyExists) {
				t.Log.Warnf("Repository already exists. Skipping clone.")

				r, err = git.PlainOpen(a.workingDirectory)
				if err != nil {
					return err
				}

				remote, err := r.Remote("origin")
				if err != nil {
					return err
				}

				if remote.Config().URLs[0] != a.config.Repository {
					return fmt.Errorf("Remote repository URL does not match the provided URL: %s -> %s", remote.Config().URLs[0], a.config.Repository)
				}
			} else if err != nil {
				return err
			}

			a.repository = r

			ref, err := r.Head()
			if err != nil {
				return err
			}
			t.Log.Infof("Repository cloned successfully: %s", ref)

			a.state = &GitAdapterState{
				LastCommit: ref.Hash().String(),
			}

			return nil
		}).
		Job()
}

func (a *GitAdapter) Sync() Job {
	return a.tl.CreateTask("sync").
		Set(func(t *Task[any]) error {
			t.Log.Infof("Pulling repository: %s@%s", a.config.Repository, a.config.Branch)

			w, err := a.repository.Worktree()
			if err != nil {
				return err
			}

			originalRef, err := a.repository.Head()
			if err != nil {
				return err
			}

			err = w.Pull(&git.PullOptions{
				Progress:      t.Log.Writer(),
				Auth:          a.authMethod,
				ReferenceName: plumbing.ReferenceName("HEAD"),
			})
			if errors.Is(err, git.NoErrAlreadyUpToDate) {
				ref, err := a.repository.Head()
				if err != nil {
					return err
				}

				t.Log.Infof("Repository is already up-to-date: %s", ref)

				return nil
			} else if err != nil {
				return err
			}

			ref, err := a.repository.Head()
			if err != nil {
				return err
			}

			a.state.LastCommit = ref.Hash().String()
			a.ctx.State.SetDirty()

			t.Log.Infof("Repository pulled successfully: %s -> %s", originalRef, ref)

			return nil
		}).
		Job()
}

func (a *GitAdapter) Finalize() Job {
	return a.tl.CreateTask("finalize").
		Set(func(t *Task[any]) error {
			t.CreateSubtask("sync", "delete").
				ShouldDisable(func(_ *Task[any]) bool {
					return !a.ctx.Flags.ForceSync && (!a.ctx.Flags.SyncDelete || !a.ctx.State.IsDirty())
				}).
				Set(func(t *Task[any]) error {
					f, err := a.ctx.State.Read()
					if err != nil {
						return err
					}
					state := &GitAdapterState{}
					if f == nil {
						t.Log.Warnf("State file does not exists.")

						return nil
					}

					if err := json.Unmarshal(f, &state); err != nil {
						t.Log.Errorf("Failed to unmarshal state file: %v", err)

						return nil
					}

					t.Log.Debugf("Syncing deleted files: from commit %s", state.LastCommit)

					last, err := a.repository.CommitObject(plumbing.NewHash(state.LastCommit))
					if err != nil {
						return err
					}
					head, err := a.repository.Head()
					if err != nil {
						return err
					}
					now, err := a.repository.CommitObject(head.Hash())
					if err != nil {
						return err
					}

					diff, err := last.Patch(now)
					if err != nil {
						return err
					}

					patches := diff.FilePatches()

					if len(patches) == 0 {
						t.Log.Infof("No changes found between %s and %s", last.Hash, now.Hash)

						return nil
					}

					for _, file := range patches {
						from, to := file.Files()
						//nolint: nestif
						if to == nil {
							path, err := filepath.Rel(a.ctx.RootDirectory, fmt.Sprintf("/%s", from.Path()))
							if err != nil {
								return err
							}

							tf := operations.NewFile(filepath.Join(a.ctx.TargetDirectory, path))

							err = tf.Remove()
							if err != nil {
								t.Log.Warnf("File already did not exists: %s", path)
							} else {
								t.Log.Warnf("File deleted: %s", path)
							}

							ls, err := tf.ReadDir()
							if err != nil {
								return err
							}
							if len(ls) == 0 && a.ctx.Flags.SyncDeleteEmptyDirectories {
								err = os.Remove(tf.Cwd())
								if err != nil {
									return err
								}

								t.Log.Warnf("Empty directory deleted: %s", tf.Cwd())
							}
						}
					}

					return nil
				}).
				AddSelfToTheParentAsSequence()

			t.CreateSubtask("state", "commit").
				Set(func(t *Task[any]) error {
					a.ctx.State.SetClean()

					f, err := json.Marshal(a.state)
					if err != nil {
						return err
					}

					t.Log.Debugf("Writing state file: %s", f)

					return a.ctx.State.Write(f)
				}).
				AddSelfToTheParentAsSequence()

			return t.RunSubtasks()
		}).
		Job()
}
