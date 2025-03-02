package internal

import (
	"github.com/sirupsen/logrus"
	"gitlab.kilic.dev/docker/beamer/internal/operations"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type State struct {
	ctx   *ServiceCtx
	file  string
	dirty bool
	log   *logrus.Entry
}

func NewState(ctx *ServiceCtx, file string) *State {
	return &State{
		ctx:  ctx,
		file: file,
		log:  ctx.Log.WithField(LOG_FIELD_CONTEXT, "state"),
	}
}

func (s *State) Read() ([]byte, error) {
	s.log.Debugf("Reading state: %s", s.file)

	f := operations.NewFile(s.file)
	if !f.Exists() {
		return nil, nil
	}

	return f.ReadFile()
}

func (s *State) Write(data []byte) error {
	s.log.Debugf("Writing state: %s", s.file)

	f := operations.NewFile(s.file)

	return f.WriteFile(data, 0600)
}

func (s *State) SetDirty() {
	s.dirty = true
}

func (s *State) IsDirty() bool {
	return s.dirty
}

func (s *State) SetClean() {
	s.dirty = false
}
