// nolint: cyclop
package pipe

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	glob "github.com/bmatcuk/doublestar/v4"
	. "gitlab.kilic.dev/libraries/plumber/v5"
	"golang.org/x/sync/errgroup"
)

// nolint: gocyclo, cyclop
func Workflow(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("workflow").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Ctx.Fetch.Dirty && !t.Pipe.Config.ForceWorkflow
		}).
		Set(func(t *Task[Pipe]) error {
			source := filepath.Join(t.Pipe.Config.WorkingDirectory, t.Pipe.Config.RootDirectory)

			t.Log.Debugf("Walking from directory: %s", source)

			ignored := []string{
				fmt.Sprintf("**/%s", t.Pipe.Config.IgnoreFile),
				".git/**",
			}

			if t.Pipe.Config.IgnoreFile != "" {
				file := filepath.Join(source, t.Pipe.Config.IgnoreFile)

				stat, err := os.Stat(file)
				if err != nil || stat.IsDir() {
					t.Log.Debugf("Ignore file not found: %s", file)
				} else {
					f, err := os.Open(file)
					if err != nil {
						return err
					}
					defer f.Close()

					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						line := scanner.Text()
						if line == "" {
							ignored = append(ignored, line)
						}
					}
				}
			}

			t.Log.Debugf("Ignoring patterns: %v", ignored)

			files := []string{}

			g := errgroup.Group{}
			err := filepath.WalkDir(
				source,
				func(abs string, d fs.DirEntry, e error) error {
					if e != nil {
						return fmt.Errorf("Error walking: %s -> %w", abs, e)
					} else if d.IsDir() {
						return nil
					}

					path, err := filepath.Rel(t.Pipe.Config.WorkingDirectory, abs)
					if err != nil {
						return err
					}

					for _, pattern := range ignored {
						match, err := glob.PathMatch(pattern, path)
						if err != nil {
							return err
						} else if match {
							t.Log.Debugf("Ignoring: %s", path)

							return nil
						}
					}

					files = append(files, path)

					return nil
				},
			)
			if err != nil {
				return err
			}

			err = g.Wait()
			if err != nil {
				return err
			}

			// create directories
			dirs := []string{}
			for _, path := range files {
				dirs = append(dirs, filepath.Dir(path))
			}
			dirs = slices.Compact(dirs)

			for _, dir := range dirs {
				g.Go(func() error {
					path := filepath.Join(t.Pipe.Config.TargetDirectory, dir)

					stat, err := os.Stat(filepath.Join(source, dir))
					if err != nil {
						return fmt.Errorf("Can not get the stat of the source directory: %s -> %w", filepath.Join(source, dir), err)
					}
					perm := stat.Mode().Perm()

					t.Log.Debugf("Directory needed in target: %s with %s in %s", dir, perm, t.Pipe.Config.TargetDirectory)

					err = os.MkdirAll(path, perm)
					if err != nil {
						return err
					}

					return nil
				})
			}

			err = g.Wait()
			if err != nil {
				return err
			}

			// process files

			for _, path := range files {
				g.Go(func() error {
					t.Log.Debugf("Processing: %s", path)

					tf := filepath.Join(t.Pipe.Config.TargetDirectory, path)
					sf := filepath.Join(source, path)

					ss, err := os.Stat(tf)
					if err == nil && ss.IsDir() {
						return fmt.Errorf("Target is a directory: %s", path)
					}
					if errors.Is(err, os.ErrNotExist) {
						t.Log.Debugf("File already does not exists copying to target: %s", tf)

						src, err := os.Open(sf)
						if err != nil {
							return err
						}
						defer src.Close()

						dest, err := os.Create(tf)
						if err != nil {
							return err
						}
						defer dest.Close()

						_, err = io.Copy(dest, src)
						if err != nil {
							return err
						}
					} else if err != nil {
						return err
					} else {
						t.Log.Debugf("File already exists: %s", tf)

						f1, err := os.Open(tf)
						if err != nil {
							return err
						}
						defer f1.Close()

						f2, err := os.Open(sf)
						if err != nil {
							return err
						}
						defer f2.Close()

						h1 := sha256.New()
						if _, err := io.Copy(h1, f1); err != nil {
							return err
						}

						h2 := sha256.New()
						if _, err := io.Copy(h2, f2); err != nil {
							return err
						}

						if bytes.Equal(h1.Sum(nil), h2.Sum(nil)) {
							t.Log.Debugf("Files are the same, nothing to do: %s -> %s", path, tf)
						} else {
							t.Log.Infof("File has changed, updating: %s -> %s", path, tf)

							src, err := os.Open(sf)
							if err != nil {
								return err
							}
							defer src.Close()

							dest, err := os.Create(tf)
							if err != nil {
								return err
							}
							defer dest.Close()

							_, err = io.Copy(dest, src)
							if err != nil {
								return err
							}
						}
					}

					ss, err = os.Stat(sf)
					if err != nil {
						return err
					}
					ts, err := os.Stat(tf)
					if err != nil {
						return err
					}
					if ss.Mode().Perm() != ts.Mode().Perm() {
						perm := ss.Mode().Perm()
						t.Log.Debugf("Setting permissions for: %s with %s", tf, perm)

						err = os.Chmod(tf, perm)
						if err != nil {
							return err
						}
					}

					return nil
				})
			}

			err = g.Wait()
			if err != nil {
				return err
			}

			f := filepath.Join(t.Pipe.Config.TargetDirectory, t.Pipe.Config.StateFile)
			t.Log.Debugf("Writing state file: %s -> %s", f, t.Pipe.Ctx.Fetch.State)
			err = os.WriteFile(f, t.Pipe.Ctx.Fetch.State, 0600)
			if err != nil {
				return err
			}

			t.Pipe.Ctx.Fetch.Dirty = false

			return nil
		})
}
