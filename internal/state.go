package services

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
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

	f, err := os.ReadFile(s.file)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return f, err
}

func (s *State) Write(data []byte) error {
	s.log.Debugf("Writing state: %s", s.file)

	return os.WriteFile(s.file, data, 0600)
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
