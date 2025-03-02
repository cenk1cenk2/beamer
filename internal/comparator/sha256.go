package comparator

import (
	"bytes"
	"crypto/sha256"
	"io"

	"gitlab.kilic.dev/docker/beamer/internal/operations"
)

type FileComparatorSha256 struct{}

var _ FileComparator = (*FileComparatorSha256)(nil)

func NewFileComparatorSha256() *FileComparatorSha256 {
	return &FileComparatorSha256{}
}

func (f *FileComparatorSha256) CompareFiles(a *operations.File, b *operations.File) (bool, error) {
	if a == nil || b == nil {
		return false, nil
	}

	f1, err := a.OpenFile()
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := b.OpenFile()
	if err != nil {
		return false, err
	}
	defer f2.Close()

	h1 := sha256.New()
	if _, err := io.Copy(h1, f1); err != nil {
		return false, err
	}

	h2 := sha256.New()
	if _, err := io.Copy(h2, f2); err != nil {
		return false, err
	}

	return bytes.Equal(h1.Sum(nil), h2.Sum(nil)), nil
}
