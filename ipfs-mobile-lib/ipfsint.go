package ipfsApi

import (
	"errors"
	"fmt"
	"os"
	"time"

	"context"

	//path "gx/ipfs/QmQa2wf1sLFKkjHCVEbna8y5qhdMjL8vtTJSAc48vZGTer/go-ipfs/path"
	//cli "gx/ipfs/QmVcLF2CgjQb5BWmYFWsDfxDjbzBfcChfdHRedxeL3dV4K/cli"

	//path "github.com/ipfs/go-ipfs/path"

	"github.com/lemonwin798/go-ipfs-mobile/ipfs-mobile-path"
	"github.com/lemonwin798/go-ipfs-mobile/mobile-log"
	//fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	//util "gx/ipfs/QmNiJuT8Ja3hMVpBHXv3Q6dwmperaQ6JjLtpMQgMCD7xvx/go-ipfs-util"
)

var (
	mshell *Shell

	ipfspath = "/storage/emulated/0/Android/data/org.golang.todo.github_com_ipfs_go_ipfs_mobile/files/"

	bootstarps []string
)

func Api_SetBoostarp(bootstarpurl string) {
	bootstarps = append(bootstarps, bootstarpurl)
}

func Api_InitNode(tmpNode bool, privateKey string) error {

	if mshell != nil {
		mobileLog.Print("InitNode repeated finish===")
		return nil
	}

	err := mobilePath.InitMobilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "InitIpfsStoragePath error: %s\n", err)
		return err
	}

	ipfspath = mobilePath.GetExternStorageFilePath()

	err = mobileLog.NewMobileLog(ipfspath + "log.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewMobileLog error: %s\n", err)
		return err
	}

	ipfspath = mobilePath.GetExternStorageFilePath() + "ipfs-shell/"
	os.Setenv("IPFS_PATH", ipfspath)

	mobileLog.Print("os Getenv succeed===", ipfspath)

	shell, err := newInternalShell(tmpNode, privateKey, bootstarps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	mshell = shell

	mobileLog.Print("InitNode finish===")
	return nil
}

func Api_CloseNode() {
	mshell.CloseShell()
	ClearNode()
	mobileLog.CloseMobileLog()
	mshell = nil
	bootstarps = []string{}
}

func Api_ServeHTTPGateway() error {

	if mshell == nil {
		return errors.New("shell is null")
	}
	//var gwErrc <-chan error

	var err error
	_, err = mshell.ServeHTTPGateway(ipfspath)
	if err != nil {
		return err
	}

	return nil
}

func Api_Get(path string, outfile string) {

	savefile := ipfspath + outfile
	mobileLog.Print("Api_Get: save:", savefile)

	begin := time.Now()

	if err := mshell.Get(path, savefile); err != nil {
		os.Remove(savefile)
		mobileLog.Print("ipget Get failed: %s\n", err)
		os.Exit(2)
	}
	d := time.Since(begin)

	mobileLog.Print("Get: gas time: %u ", d, "save:", savefile)
}

func Api_Catching(path string) []byte {

	begin := time.Now()

	buf, err := mshell.Catching(path)
	if err != nil {

		mobileLog.Print("ipget Catching failed: %s\n", err)
		os.Exit(2)
	}
	d := time.Since(begin)

	mobileLog.Print("Catching: gas time: %u", d)
	return buf
}

func newInternalShell(tmpNode bool, privateKey string, bootstarpurl []string) (*Shell, error) {
	ctx, _ := context.WithCancel(context.Background())

	// Cancel the ipfs node context if the process gets interrupted or killed.
	// TODO(noffle): is this needed?
	//go func() {
	//	interrupts := make(chan os.Signal, 1)
	//	signal.Notify(interrupts, os.Interrupt, os.Kill)
	//	<-interrupts
	//	cancel()
	//}()

	/*shell, err := tryLocal(ctx)
	if err == nil {
		return shell, nil
	}*/

	mobileLog.Print("mobileShell.NewMobileNode---")

	node, err := NewMobileNode(ctx, ipfspath, privateKey, bootstarpurl, tmpNode)
	if err != nil {
		mobileLog.Print("NewMobileNode: %s", err)
		return nil, err
	}
	return NewShell(node), nil
}
