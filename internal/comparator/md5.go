// nolint: gosec, dupl
package comparator

import (
	"bytes"
	"crypto/md5"
	"io"

	"gitlab.kilic.dev/docker/beamer/internal/operations"
)

type FileComparatorMd5 struct{}

var _ FileComparator = (*FileComparatorMd5)(nil)

func NewFileComparatorMd5() *FileComparatorMd5 {
	return &FileComparatorMd5{}
}

func (f *FileComparatorMd5) CompareFiles(a *operations.File, b *operations.File) (bool, error) {
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

	h1 := md5.New()
	if _, err := io.Copy(h1, f1); err != nil {
		return false, err
	}

	h2 := md5.New()
	if _, err := io.Copy(h2, f2); err != nil {
		return false, err
	}

	return bytes.Equal(h1.Sum(nil), h2.Sum(nil)), nil
}
