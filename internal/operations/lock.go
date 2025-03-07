package operations

type LockFile struct {
	file *File
}

func NewLockFile(path ...string) *LockFile {
	file := NewFile(path...)

	return &LockFile{
		file: file,
	}
}

func (f *LockFile) Lock() error {
	return f.file.Touch()
}

func (f *LockFile) Unlock() error {
	return f.file.Remove()
}

func (f *LockFile) IsLocked() bool {
	return f.file.Exists()
}

func (f *LockFile) Path() string {
	return f.file.Abs()
}
