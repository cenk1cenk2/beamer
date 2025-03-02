package pipe

import (
	"gitlab.kilic.dev/docker/beamer/internal/comparator"
)

type Ctx struct {
	FileComparator comparator.FileComparator

	Fetch struct {
		Dirty bool
		State []byte
	}
}
