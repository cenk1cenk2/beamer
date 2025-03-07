package pipe

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	glob "github.com/bmatcuk/doublestar/v4"
	"gitlab.kilic.dev/docker/beamer/internal/operations"
	. "gitlab.kilic.dev/libraries/plumber/v5"
	"golang.org/x/sync/errgroup"
)

func Workflow(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("workflow").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Ctx.State.IsDirty() && !t.Pipe.Config.ForceWorkflow
		}).
		Set(func(t *Task[Pipe]) error {
			ignored, err := parseIgnoreFile(t)
			if err != nil {
				return err
			}
			t.Log.Debugf("Ignoring patterns: %v", ignored)

			files, err := walkdir(t, ignored)
			if err != nil {
				return err
			}
			t.Log.Debugf("Files to process: %v", files)

			// create directories
			if err := ensureDirs(t, files); err != nil {
				return err
			}

			// process files

			g := errgroup.Group{}
			for _, path := range files {
				g.Go(func() error {
					return processFile(t, path)
				})
			}

			err = g.Wait()
			if err != nil {
				return err
			}

			return nil
		})
}

func parseIgnoreFile(t *Task[Pipe]) ([]string, error) {
	ignored := []string{
		fmt.Sprintf("**/%s", t.Pipe.Config.IgnoreFile),
		".git/**",
	}

	if t.Pipe.Config.IgnoreFile == "" {
		return ignored, nil
	}

	f := operations.NewFile(t.Pipe.WorkingDirectory, t.Pipe.RootDirectory, t.Pipe.Config.IgnoreFile)

	if !f.IsFile() {
		t.Log.Debugf("Ignore file not found: %s", f.Abs())

		return ignored, nil
	}

	lines, err := f.ReadLines()
	if err != nil {
		return nil, err
	}

	ignored = append(ignored, lines...)

	return ignored, nil
}

func walkdir(t *Task[Pipe], ignored []string) ([]string, error) {
	files := []string{}

	f := operations.NewFile(t.Pipe.WorkingDirectory, t.Pipe.RootDirectory)

	t.Log.Debugf("Walking from source directory: %s", f.Abs())

	g := errgroup.Group{}
	err := filepath.WalkDir(
		f.Abs(),
		func(abs string, d fs.DirEntry, e error) error {
			if e != nil {
				return fmt.Errorf("Error walking: %s -> %w", abs, e)
			} else if d.IsDir() {
				return nil
			}

			f := operations.NewFile(abs)

			path, err := f.RelTo(t.Pipe.WorkingDirectory)
			if err != nil {
				return err
			}

			for _, pattern := range ignored {
				match, err := glob.PathMatch(pattern, path)
				if err != nil {
					return err
				} else if match {
					t.Log.Debugf("Ignoring: %s with /%s/", path, pattern)

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

func ensureDirs(t *Task[Pipe], files []string) error {
	g := errgroup.Group{}

	dirs := []string{}

	for _, path := range files {
		dirs = append(dirs, filepath.Dir(path))
	}
	dirs = slices.Compact(dirs)

	for _, dir := range dirs {
		g.Go(func() error {
			source := operations.NewFile(t.Pipe.WorkingDirectory, dir)
			rel, err := filepath.Rel(t.Pipe.RootDirectory, fmt.Sprintf("/%s", dir))
			if err != nil {
				return err
			}
			target := operations.NewFile(t.Pipe.TargetDirectory, rel)

			if !source.IsDir() {
				return fmt.Errorf("Source is not a directory anymore: %s", source.Abs())
			} else if target.IsDir() {
				t.Log.Debugf("Directory already exists in target: %s", target.Abs())

				return nil
			}

			stat, err := source.Stat()
			if err != nil {
				return err
			}

			t.Log.Debugf("Directory needed in target: %s with %s in %s", target.Rel(), stat.Mode().Perm(), target.Cwd())

			return target.Mkdirp(stat.Mode())
		})
	}

	return g.Wait()
}

func processFile(t *Task[Pipe], path string) error {
	sf := operations.NewFile(t.Pipe.WorkingDirectory, path)
	rel, err := filepath.Rel(t.Pipe.RootDirectory, fmt.Sprintf("/%s", path))
	if err != nil {
		return err
	}
	tf := operations.NewFile(t.Pipe.TargetDirectory, rel)

	t.Log.Debugf("Processing: %s -> %s", sf.Abs(), tf.Abs())

	if sf.IsDir() {
		return fmt.Errorf("Source is a directory: %s", sf.Abs())
	} else if tf.IsDir() {
		return fmt.Errorf("Target is a directory: %s", tf.Abs())
	}

	//nolint:nestif
	if slices.Contains(t.Pipe.TemplateFiles, tf.Ext()) {
		f, err := sf.ReadFile()
		if err != nil {
			return err
		}

		template, err := InlineTemplate(string(f), struct{}{})
		if err != nil {
			return fmt.Errorf("Can not inline template for file: %s -> %w", sf.Abs(), err)
		}

		ss, err := sf.Stat()
		if err != nil {
			return err
		}

		// change the source file to the templated file
		path := strings.TrimSuffix(path, sf.Ext())

		temp, err := os.CreateTemp("", fmt.Sprintf("beamer-%s", filepath.Base(path)))
		if err != nil {
			return err
		}
		defer os.Remove(temp.Name())

		nf := operations.NewFile(temp.Name())
		if err := nf.WriteFile([]byte(template), ss.Mode().Perm()); err != nil {
			return err
		}

		tf = operations.NewFile(t.Pipe.TargetDirectory, path)

		t.Log.Debugf("Templated file: %s (from temp %s) -> %s", sf.Abs(), nf.Abs(), tf.Abs())

		sf = nf
	}

	if tf.Exists() {
		t.Log.Debugf("File already exists: %s", tf.Abs())

		equal, err := t.Pipe.Ctx.FileComparator.CompareFiles(sf, tf)
		if err != nil {
			return err
		}

		if equal {
			t.Log.Debugf("Files are the same, nothing to do: %s -> %s", sf.Abs(), tf.Abs())

			return nil
		}

		t.Log.Infof("File has changed, updating: %s -> %s", sf.Abs(), tf.Abs())
	} else {
		t.Log.Debugf("File already does not exists copying to target: %s", tf.Abs())
	}

	return sf.CopyTo(tf)
}
