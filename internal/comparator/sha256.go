package comparator

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
)

type FileComparatorSha256 struct{}

var _ FileComparator = (*FileComparatorSha256)(nil)

func NewFileComparatorSha256() *FileComparatorSha256 {
	return &FileComparatorSha256{}
}

func (f *FileComparatorSha256) CompareFiles(f1 *os.File, f2 *os.File) (bool, error) {
	if f1 == nil || f2 == nil {
		return false, nil
	}

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
