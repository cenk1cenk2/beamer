package services

import (
	"github.com/sirupsen/logrus"
)

type ServiceCtx struct {
	Log              *logrus.Entry
	State            *State
	WorkingDirectory string
	RootDirectory    string
	TargetDirectory  string
	Flags            *ServiceFlags
}

type ServiceFlags struct {
	WorkingDirectory string
	TargetDirectory  string `validate:"required"`
	RootDirectory    string

	ForceSync                  bool
	SyncDelete                 bool
	SyncDeleteEmptyDirectories bool
}
