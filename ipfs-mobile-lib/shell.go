package ipfsApi

import (
	"fmt"
	"io"
	//"sync"
	"mime"

	"gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"

	oldcmds "github.com/ipfs/go-ipfs/commands"
	corehttp "github.com/ipfs/go-ipfs/core/corehttp"
	config "github.com/ipfs/go-ipfs/repo/config"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"

	"github.com/ipfs/go-ipfs/core"
	"github.com/lemonwin798/go-ipfs-mobile/mobile-log"
	//"golang.org/x/net/context"
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
	//ctx         context.Context
	node        *core.IpfsNode
	chanGateway chan error
	gwLis       manet.Listener
}

// func NewReadOnlyShell() *Shell {}

func NewShell(node *core.IpfsNode) *Shell {
	//return NewShellWithContext(node, context.Background())
	return &Shell{node, nil, nil}
}

//func NewShellWithContext(node *core.IpfsNode, ctx context.Context) *Shell {
//	return &Shell{ctx, node, nil, nil}
//}

func (s *Shell) ServeHTTPGateway(repoPath string) (<-chan error, error) {
	mobileLog.Print("ServeHTTPGateway begin===\n")
	if s.gwLis != nil {
		mobileLog.Print("ServeHTTPGateway repeated===\n")
		return nil, nil
	}

	cfg, err := s.node.Repo.Config()
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: GetConfig() failed: %s", err)
	}

	gatewayMaddr, err := ma.NewMultiaddr(cfg.Addresses.Gateway)
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", cfg.Addresses.Gateway, err)
	}

	writable := false //req.Options[writableKwd].(bool)
	//if !writableOptionFound {
	//	writable = cfg.Gateway.Writable
	//}

	s.gwLis, err = manet.Listen(gatewayMaddr)
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
	}
	mobileLog.Print("manet.Listen===" + cfg.Addresses.Gateway + "\n")
	// we might have listened to /tcp/0 - lets see what we are listing on
	gatewayMaddr = s.gwLis.Multiaddr()

	if writable {
		fmt.Printf("Gateway (writable) server listening on %s\n", gatewayMaddr)
	} else {
		fmt.Printf("Gateway (readonly) server listening on %s\n", gatewayMaddr)
	}

	cctx := &oldcmds.Context{
		ConfigRoot: repoPath,
		LoadConfig: loadConfig,
		ReqLog:     &oldcmds.ReqLog{},
		ConstructNode: func() (n *core.IpfsNode, err error) {
			//s.node.SetLocal(true)
			return s.node, nil
		},
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("gateway"),
		corehttp.CheckVersionOption(),
		corehttp.CommandsROOption(*cctx),
		corehttp.VersionOption(),
		corehttp.IPNSHostnameOption(),
		corehttp.GatewayOption(writable, "/ipfs", "/ipns"),
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	//在播放ts流的时候，如果不设置数据类型，则fs.go的serveContent函数每次都会把ts在内存中读取出来，进行类型匹配，
	//这在手机上会造成大量的性能消耗，从而引发客户端连接异常断开，具体原因不明
	mime.AddExtensionType(".ts", "application/octet-stream")
	mime.AddExtensionType(".m3u8", "text/plain; charset=utf-8")

	chanGateway := make(chan error)
	go func() {
		chanGateway <- corehttp.Serve(s.node, s.gwLis.NetListener(), opts...)
		close(chanGateway)

		mobileLog.Print("corehttp.Serve===\n")
	}()
	mobileLog.Print("ServeHTTPGateway finish===\n")
	return chanGateway, nil
}

func loadConfig(path string) (*config.Config, error) {
	return fsrepo.ConfigAt(path)
}

func (s *Shell) CloseShell() {
	if s.chanGateway != nil {
		mobileLog.Print("ServeHTTPGateway close===\n")
		close(s.chanGateway)
	}
	if s.gwLis != nil {
		s.gwLis.Close()
	}
	s.chanGateway = nil
	s.gwLis = nil
}

/*
func (s *Shell) Merge(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	for _, c := range cs {
		if c != nil {
			wg.Add(1)
			go output(c)
		}
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}*/
