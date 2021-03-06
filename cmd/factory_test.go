package cmd

import (
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/qri-io/qri/config"
	"github.com/qri-io/qri/lib"
	"github.com/qri-io/qri/p2p"
	"github.com/qri-io/qri/repo"
	"github.com/qri-io/qri/repo/test"
	"github.com/qri-io/registry/regclient"
)

// TestFactory is an implementation of the Factory interface for testing purposes
type TestFactory struct {
	IOStreams
	// QriRepoPath is the path to the QRI repository
	qriRepoPath string
	// IpfsFsPath is the path to the IPFS repo
	ipfsFsPath string

	// Configuration object
	config *config.Config
	node   *p2p.QriNode
	repo   repo.Repo
	rpc    *rpc.Client
}

// NewTestFactory creates TestFactory object with an in memory test repo
// with an optional registry client. In tests users can create mock registry
// servers and pass in a client connected to that mock, or omit the registry
// client entirely for testing without a designated registry
func NewTestFactory(c *regclient.Client) (tf TestFactory, err error) {
	repo, err := test.NewTestRepo(c)
	if err != nil {
		return
	}

	cfg := config.DefaultConfigForTesting()
	tnode, err := p2p.NewTestableQriNode(repo, cfg.P2P)
	if err != nil {
		return
	}

	return TestFactory{
		qriRepoPath: "",
		ipfsFsPath:  "",

		repo:   repo,
		rpc:    nil,
		config: cfg,
		node:   tnode.(*p2p.QriNode),
	}, nil
}

// Config returns from internal state
func (t TestFactory) Config() (*config.Config, error) {
	return t.config, nil
}

// IpfsFsPath returns from internal state
func (t TestFactory) IpfsFsPath() string {
	return t.ipfsFsPath
}

// QriRepoPath returns from internal state
func (t TestFactory) QriRepoPath() string {
	return t.qriRepoPath
}

// Repo returns from internal state
func (t TestFactory) Repo() (repo.Repo, error) {
	return t.repo, nil
}

// Node returns the internal qri node from state
func (t TestFactory) Node() (*p2p.QriNode, error) {
	return t.node, nil
}

// RPC returns from internal state
func (t TestFactory) RPC() *rpc.Client {
	return nil
}

// DatasetRequests generates a lib.DatasetRequests from internal state
func (t TestFactory) DatasetRequests() (*lib.DatasetRequests, error) {
	return lib.NewDatasetRequests(t.node, t.rpc), nil
}

// RegistryRequests generates a lib.RegistryRequests from internal state
func (t TestFactory) RegistryRequests() (*lib.RegistryRequests, error) {
	return lib.NewRegistryRequests(t.node, t.rpc), nil
}

// LogRequests generates a lib.LogRequests from internal state
func (t TestFactory) LogRequests() (*lib.LogRequests, error) {
	return lib.NewLogRequests(t.node, t.rpc), nil
}

// PeerRequests generates a lib.PeerRequests from internal state
func (t TestFactory) PeerRequests() (*lib.PeerRequests, error) {
	return lib.NewPeerRequests(t.node, t.rpc), nil
}

// ProfileRequests generates a lib.ProfileRequests from internal state
func (t TestFactory) ProfileRequests() (*lib.ProfileRequests, error) {
	return lib.NewProfileRequests(t.node, t.rpc), nil
}

// SelectionRequests creates a lib.SelectionRequests from internal state
func (t TestFactory) SelectionRequests() (*lib.SelectionRequests, error) {
	return lib.NewSelectionRequests(t.repo, t.rpc), nil
}

// SearchRequests generates a lib.SearchRequests from internal state
func (t TestFactory) SearchRequests() (*lib.SearchRequests, error) {
	return lib.NewSearchRequests(t.node, t.rpc), nil
}

// RenderRequests generates a lib.RenderRequests from internal state
func (t TestFactory) RenderRequests() (*lib.RenderRequests, error) {
	return lib.NewRenderRequests(t.repo, t.rpc), nil
}

func TestEnvPathFactory(t *testing.T) {
	//Needed to clean up changes after the test has finished running
	prevQRIPath := os.Getenv("QRI_PATH")
	prevIPFSPath := os.Getenv("IPFS_PATH")

	defer func() {
		os.Setenv("QRI_PATH", prevQRIPath)
		os.Setenv("IPFS_PATH", prevIPFSPath)
	}()

	//Test variables
	emptyPath := ""
	fakePath := "fake_path"
	home, err := homedir.Dir()
	if err != nil {
		t.Fatalf("Failed to find the home directory: %s", err.Error())
	}

	tests := []struct {
		qriPath    string
		ipfsPath   string
		qriAnswer  string
		ipfsAnswer string
	}{
		{emptyPath, emptyPath, filepath.Join(home, ".qri"), filepath.Join(home, ".ipfs")},
		{emptyPath, fakePath, filepath.Join(home, ".qri"), fakePath},
		{fakePath, emptyPath, fakePath, filepath.Join(home, ".ipfs")},
		{fakePath, fakePath, fakePath, fakePath},
	}

	for i, test := range tests {
		err := os.Setenv("QRI_PATH", test.qriPath)
		if err != nil {
			t.Errorf("case %d failed to set up QRI_PATH: %s", i, err.Error())
		}

		err = os.Setenv("IPFS_PATH", test.ipfsPath)
		if err != nil {
			t.Errorf("case %d failed to set up IPFS_PATH: %s", i, err.Error())
		}

		qriResult, ipfsResult := EnvPathFactory()

		if !strings.EqualFold(qriResult, test.qriAnswer) {
			t.Errorf("case %d expected qri path to be %s, but got %s", i, test.qriAnswer, qriResult)
		}

		if !strings.EqualFold(ipfsResult, test.ipfsAnswer) {
			t.Errorf("case %d Expected ipfs path to be %s, but got %s", i, test.ipfsAnswer, ipfsResult)
		}

	}
}
