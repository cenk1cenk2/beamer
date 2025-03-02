package pipe

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type Ctx struct {
	Fetch struct {
		Dirty bool
		State []byte
	}

	Git struct {
		AuthMethod    transport.AuthMethod
		SshPrivateKey []byte
		Repository    *git.Repository
	}
}
