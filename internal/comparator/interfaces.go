package comparator

import (
	"gitlab.kilic.dev/docker/beamer/internal/operations"
)

type FileComparator interface {
	CompareFiles(a *operations.File, b *operations.File) (bool, error)
}
