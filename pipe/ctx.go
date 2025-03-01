package pipe

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type Ctx struct {
	Git struct {
		AuthMethod    transport.AuthMethod
		SshPrivateKey []byte
	}
}
