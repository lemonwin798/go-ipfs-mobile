package ipfsApi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	//util "gx/ipfs/QmNiJuT8Ja3hMVpBHXv3Q6dwmperaQ6JjLtpMQgMCD7xvx/go-ipfs-util"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo/config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/lemonwin798/go-ipfs-mobile/mobile-log"
	"golang.org/x/net/context"
)

var (
	tmpNodeDir string
	tmpNode    = false

	node *core.IpfsNode
)

func NewMobileNode(ctx context.Context, repoPath string, privateKey string, bootstarpurl []string, cacheNode bool) (*core.IpfsNode, error) {

	node, err := loadMobileNode(ctx, repoPath, privateKey, bootstarpurl)
	if err != nil {

		return nil, err
	}
	tmpNodeDir = repoPath

	tmpNode = cacheNode

	return node, nil
}

func loadMobileNode(ctx context.Context, repoPath string, privateKey string, bootstarpurl []string) (*core.IpfsNode, error) {

	r, err := fsrepo.Open(repoPath)
	if err != nil {
		//mobileLog.Print("opening fsrepo failed: %s", err)
		//return nil, fmt.Errorf("opening fsrepo failed: %s", err)
		createMobileNode(repoPath, privateKey, bootstarpurl)

		r, err = fsrepo.Open(repoPath)
		if err != nil {
			mobileLog.Print("opening fsrepo failed: %s", err)
			return nil, fmt.Errorf("opening fsrepo failed: %s", err)
		}
	}

	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online:    true,
		Repo:      r,
		Permanent: false,
		Routing:   core.DHTClientOption,
	})
	if err != nil {
		mobileLog.Print("ipfs NewNode() failed: %s", err)
		return nil, fmt.Errorf("ipfs NewNode() failed: %s", err)
	}
	node.SetLocal(false)

	// TODO: can we bootsrap localy/mdns first and fall back to default?
	err = node.Bootstrap(core.DefaultBootstrapConfig)
	if err != nil {
		mobileLog.Print("ipfs Bootstrap() failed: %s", err)
		return nil, fmt.Errorf("ipfs Bootstrap() failed: %s", err)
	}

	tmpNode = false

	return node, nil
}

func createMobileNode(repoPath string, privateKey string, bootstarpurl []string) error {
	//dir, err := ioutil.TempDir(repoPath, "ipfs-shell")
	//if err != nil {
	//	mobileLog.Print("failed to get temp dir: %s", err)
	//	return nil, fmt.Errorf("failed to get temp dir: %s", err)
	//}

	cfg, err := config.Init(ioutil.Discard, 1024)
	if err != nil {
		mobileLog.Print("config.Init(ioutil.Discard, 1024): %s", err)
		return err
	}

	if len(bootstarpurl) > 0 {

		bootstrapPeers, err := privateBootstrapPeers(bootstarpurl)
		if err != nil {
			return err
		}

		cfg.Bootstrap = config.BootstrapPeerStrings(bootstrapPeers)
	}

	cfg.Addresses.Gateway = "/ip4/127.0.0.1/tcp/8089"

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		mobileLog.Print("failed to init ephemeral node: %s", err)
		return fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	if len(privateKey) == 64 {
		keyfile := ipfspath + "/swarm.key"

		{

			strkey := "/key/swarm/psk/1.0.0/" + "\n"
			strkey += "/base16/" + "\n"

			strkey += privateKey

			err = ioutil.WriteFile(keyfile, []byte(strkey), 0666)
			if err != nil {
				mobileLog.Print("ioutil.WriteFile:[%s] %s", keyfile, err)
				return err
			}
		}
	}

	/*repo, err := fsrepo.Open(repoPath)
	if err != nil {
		mobileLog.Print("fsrepo.Open(dir): %s", err)
		return nil, err
	}

	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   repo,
	})*/
	return nil
}

func privateBootstrapPeers(bootstarpurl []string) ([]config.BootstrapPeer, error) {
	ps, err := config.ParseBootstrapPeers(bootstarpurl)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse hardcoded bootstrap peers: %s
This is a problem with the ipfs codebase. Please report it to the dev team.`, err)
	}
	return ps, nil
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
	if node != nil {
		node.Close()
		node = nil
	}

	if tmpNode == true {
		destorytemp(tmpNodeDir)
	}
}
