package adapter

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type Adapter interface {
	Init() Job
	Sync() Job
	Finalize() Job
}
