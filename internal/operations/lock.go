package operations

import "github.com/sirupsen/logrus"

type LockFile struct {
	file *File
	log  *logrus.Entry
}

func NewLockFile(log *logrus.Entry, path ...string) *LockFile {
	file := NewFile(path...)

	return &LockFile{
		log:  log,
		file: file,
	}
}

func (f *LockFile) Lock() error {
	f.log.Debugf("Locking: %s", f.file.Abs())

	return f.file.Touch()
}

func (f *LockFile) Unlock() error {
	f.log.Debugf("Unlocking: %s", f.file.Abs())

	return f.file.Remove()
}

func (f *LockFile) IsLocked() bool {
	locked := f.file.Exists()

	f.log.Debugf("IsLocked: %t", locked)

	return locked
}

func (f *LockFile) Path() string {
	return f.file.Abs()
}
