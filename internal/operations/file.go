package operations

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	cwd  string
	path string
}

func NewFile(path ...string) *File {
	p := filepath.Join(path...)

	return &File{
		cwd:  filepath.Dir(p),
		path: filepath.Base(p),
	}
}

func (f *File) Abs() string {
	return filepath.Join(f.cwd, f.path)
}

func (f *File) Rel() string {
	return f.path
}

func (f *File) Cwd() string {
	return f.cwd
}

func (f *File) Ext() string {
	return filepath.Ext(f.path)
}

func (f *File) RelTo(base string) (string, error) {
	return filepath.Rel(base, f.Abs())
}

func (f *File) Stat() (os.FileInfo, error) {
	return os.Stat(f.Abs())
}

func (f *File) Exists() bool {
	_, err := f.Stat()

	return err == nil
}

func (f *File) IsFile() bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

func (f *File) IsDir() bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func (f *File) ReadLines() ([]string, error) {
	var lines []string

	h, err := os.Open(f.Abs())
	if err != nil {
		return nil, err
	}
	defer h.Close()

	scanner := bufio.NewScanner(h)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func (f *File) OpenFile() (*os.File, error) {
	return os.Open(f.Abs())
}

func (f *File) ReadFile() ([]byte, error) {
	return os.ReadFile(f.Abs())
}

func (f *File) ReadDir() ([]os.DirEntry, error) {
	if !f.IsDir() {
		return os.ReadDir(f.Cwd())
	}

	return os.ReadDir(f.Abs())
}

func (f *File) WriteFile(data []byte, perm os.FileMode) error {
	return os.WriteFile(f.Abs(), data, perm)
}

func (f *File) Touch() error {
	_, err := os.Create(f.Abs())

	return err
}

func (f *File) MatchModeWith(target *File) error {
	ts, err := target.Stat()
	if err != nil {
		return err
	}
	ss, err := f.Stat()
	if err != nil {
		return err
	}

	if ts.Mode() == ss.Mode() {
		return nil
	}

	return f.Chmod(ts.Mode())
}

func (f *File) Chmod(perm os.FileMode) error {
	return os.Chmod(f.Abs(), perm)
}

func (f *File) CopyTo(target *File) error {
	src, err := os.Open(f.Abs())
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(target.Abs())
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return target.MatchModeWith(f)
}

func (f *File) Mkdirp(perm os.FileMode) error {
	return os.MkdirAll(f.Abs(), perm)
}

func (f *File) Remove() error {
	return os.Remove(f.Abs())
}
