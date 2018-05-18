package ipfsApi

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ipfs/go-ipfs/core"
	dag "github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/path"
	tar "github.com/ipfs/go-ipfs/thirdparty/tar"
	uarchive "github.com/ipfs/go-ipfs/unixfs/archive"
	unixfsio "github.com/ipfs/go-ipfs/unixfs/io"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Get(ref, outdir string) error {
	ipfsPath, err := path.ParsePath(ref)
	if err != nil {
		return fmt.Errorf("get: could not parse %q: %s", ref, err)
	}

	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, ipfsPath)
	if err != nil {
		return fmt.Errorf("get: could not resolve %s: %s", ipfsPath, err)
	}

	pbnd, ok := nd.(*dag.ProtoNode)
	if !ok {
		return errors.New("could not cast Node to ProtoNode")
	}

	r, err := uarchive.DagArchive(s.ctx, pbnd, outdir, s.node.DAG, false, 0)
	if err != nil {
		return err
	}

	ext := tar.Extractor{outdir, ProgressExtract}

	return ext.Extract(r)
}

func ProgressExtract(v int64) int64 {
	return v
}

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Cat(p string) (io.ReadCloser, error) {
	ipfsPath, err := path.ParsePath(p)
	if err != nil {
		return nil, fmt.Errorf("cat: could not parse %q: %s", p, err)
	}
	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, ipfsPath)
	if err != nil {
		return nil, fmt.Errorf("cat: could not resolve %s: %s", ipfsPath, err)
	}
	dr, err := unixfsio.NewDagReader(s.ctx, nd, s.node.DAG)
	if err != nil {
		return nil, fmt.Errorf("cat: failed to construct DAG reader: %s", err)
	}

	return dr, nil
}

func (s *Shell) Catching(p string) ([]byte, error) {
	rc, err := s.Cat(p)
	if err != nil {
		return nil, err
	}
	//fmt.Print("======size :" + strconv.Itoa(int(rc.Size())) + "\n")

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	//ioutil.WriteFile("test.gif", buf, 0666)

	return buf, nil
}
