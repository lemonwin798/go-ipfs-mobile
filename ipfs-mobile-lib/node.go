package ipfsApi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/lemonwin798/go-ipfs-mobile/mobile-log"
	"golang.org/x/net/context"
)

var (
	tmpNodeDir string
	tmpNode    = false
)

func NewMobileNode(ctx context.Context, repoPath string, cacheNode bool) (*core.IpfsNode, error) {

	dir := repoPath + "ipfs-shell"

	node, err := LoadMobileNode(ctx, dir)
	if err != nil {
		node, err = CreateMobileNode(ctx, dir)
		if err != nil {
			return nil, err
		}
	}

	tmpNodeDir = dir

	tmpNode = cacheNode
	return node, nil
}

func LoadMobileNode(ctx context.Context, repoPath string) (*core.IpfsNode, error) {

	r, err := fsrepo.Open(repoPath)
	if err != nil {
		mobileLog.Print("opening fsrepo failed: %s", err)
		return nil, fmt.Errorf("opening fsrepo failed: %s", err)
	}

	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   r,
	})
	if err != nil {
		mobileLog.Print("ipfs NewNode() failed: %s", err)
		return nil, fmt.Errorf("ipfs NewNode() failed: %s", err)
	}

	// TODO: can we bootsrap localy/mdns first and fall back to default?
	err = node.Bootstrap(core.DefaultBootstrapConfig)
	if err != nil {
		mobileLog.Print("ipfs Bootstrap() failed: %s", err)
		return nil, fmt.Errorf("ipfs Bootstrap() failed: %s", err)
	}

	tmpNode = false

	return node, nil
}

func CreateMobileNode(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	//dir, err := ioutil.TempDir(repoPath, "ipfs-shell")
	//if err != nil {
	//	mobileLog.Print("failed to get temp dir: %s", err)
	//	return nil, fmt.Errorf("failed to get temp dir: %s", err)
	//}

	cfg, err := config.Init(ioutil.Discard, 1024)
	if err != nil {
		mobileLog.Print("config.Init(ioutil.Discard, 1024): %s", err)
		return nil, err
	}

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		mobileLog.Print("failed to init ephemeral node: %s", err)
		return nil, fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		mobileLog.Print("fsrepo.Open(dir): %s", err)
		return nil, err
	}

	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   repo,
	})
}

func destorytemp(path string) {

	filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		name := fi.Name()

		if strings.Contains(name, "ipfs-shell") {

			fmt.Println("temp file name:", path)

			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("delet dir error:", err)
			}
		}
		return nil
	})

}
func ClearNode() {
	if tmpNode == true {
		destorytemp(tmpNodeDir)
	}
}
