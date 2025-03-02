package pipe

import (
	"bufio"
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

func Workflow(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("workflow").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Ctx.Fetch.Dirty && !t.Pipe.Config.ForceWorkflow
		}).
		Set(func(t *Task[Pipe]) error {
			source := filepath.Join(t.Pipe.WorkingDirectory, t.Pipe.RootDirectory)

			t.Log.Debugf("Walking from directory: %s", source)

			ignored, err := parseIgnoreFile(t, filepath.Join(source, t.Pipe.Config.IgnoreFile))
			if err != nil {
				return err
			}
			t.Log.Debugf("Ignoring patterns: %v", ignored)

			files, err := walkdir(t, source, ignored)
			if err != nil {
				return err
			}
			t.Log.Debugf("Files to process: %v", files)

			// create directories
			if err := ensureDirs(t, source, files); err != nil {
				return err
			}

			// process files

			g := errgroup.Group{}
			for _, path := range files {
				g.Go(func() error {
					return processFile(t, source, path)
				})
			}

			err = g.Wait()
			if err != nil {
				return err
			}

			return nil
		})
}

func parseIgnoreFile(t *Task[Pipe], file string) ([]string, error) {
	ignored := []string{
		fmt.Sprintf("**/%s", t.Pipe.Config.IgnoreFile),
		".git/**",
	}

	if file == "" {
		return []string{}, nil
	}

	stat, err := os.Stat(file)

	if err != nil || stat.IsDir() {
		t.Log.Debugf("Ignore file not found: %s", file)

		return ignored, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			ignored = append(ignored, line)
		}
	}

	return ignored, nil
}

func walkdir(t *Task[Pipe], source string, ignored []string) ([]string, error) {
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

			path, err := filepath.Rel(t.Pipe.WorkingDirectory, abs)
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
		return nil, err
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return files, nil
}

func ensureDirs(t *Task[Pipe], source string, files []string) error {
	g := errgroup.Group{}

	dirs := []string{}

	for _, path := range files {
		dirs = append(dirs, filepath.Dir(path))
	}
	dirs = slices.Compact(dirs)

	for _, dir := range dirs {
		g.Go(func() error {
			path := filepath.Join(t.Pipe.TargetDirectory, dir)

			stat, err := os.Stat(filepath.Join(source, dir))
			if err != nil {
				return fmt.Errorf("Can not get the stat of the source directory: %s -> %w", filepath.Join(source, dir), err)
			}
			perm := stat.Mode().Perm()

			t.Log.Debugf("Directory needed in target: %s with %s in %s", dir, perm, t.Pipe.TargetDirectory)

			err = os.MkdirAll(path, perm)
			if err != nil {
				return err
			}

			return nil
		})
	}

	return g.Wait()
}

func processFile(t *Task[Pipe], source string, path string) error {
	t.Log.Debugf("Processing: %s", path)

	tf := filepath.Join(t.Pipe.TargetDirectory, path)
	sf := filepath.Join(source, path)

	ss, err := os.Stat(tf)
	if err == nil && ss.IsDir() {
		return fmt.Errorf("Target is a directory: %s", path)
	}

	// nolint: nestif
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

		equal, err := t.Pipe.Ctx.FileComparator.CompareFiles(f1, f2)
		if err != nil {
			return err
		}

		if equal {
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
}
