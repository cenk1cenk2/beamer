package pipe

import (
	"gitlab.kilic.dev/docker/beamer/internal"
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
)

type Ctx struct {
	FileComparator comparator.FileComparator
	State          *internal.State
}
