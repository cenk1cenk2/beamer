package pipe

import (
	"gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
	"gitlab.kilic.dev/docker/beamer/internal/operations"
)

type Ctx struct {
	FileComparator comparator.FileComparator
	State          *internal.State
	LockFile       *operations.LockFile
}
