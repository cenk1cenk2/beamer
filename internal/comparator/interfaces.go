package comparator

import (
	"os"
)

type FileComparator interface {
	CompareFiles(f1 *os.File, f2 *os.File) (bool, error)
}
