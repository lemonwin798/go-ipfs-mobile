package ipfsApi

import (
	"io"

	"github.com/ipfs/go-ipfs/core"
	"golang.org/x/net/context"
)

// Interface ...
type Interface interface {

	// Cat returns a reader returning the data under the IPFS path
	Cat(path string) (io.ReadCloser, error)

	Get(hash string, outdir string) error

	Catching(p string) ([]byte, error)
}

// Shell ...
type Shell struct {
	ctx  context.Context
	node *core.IpfsNode
}

// func NewReadOnlyShell() *Shell {}

func NewShell(node *core.IpfsNode) *Shell {
	return NewShellWithContext(node, context.Background())
}

func NewShellWithContext(node *core.IpfsNode, ctx context.Context) *Shell {
	return &Shell{ctx, node}
}
