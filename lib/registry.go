package lib

import (
	"fmt"
	"net/rpc"

	"github.com/qri-io/qri/config"
	"github.com/qri-io/qri/p2p"
	"github.com/qri-io/qri/repo"
	"github.com/qri-io/qri/repo/actions"
	"github.com/qri-io/registry"
)

// RegistryRequests defines business logic for working with registries
type RegistryRequests struct {
	node *p2p.QriNode
	repo actions.Registry
	cli  *rpc.Client
}

// CoreRequestsName implements the Requests interface
func (RegistryRequests) CoreRequestsName() string { return "registry" }

// NewRegistryRequests creates a DatasetRequests pointer from either a repo
// or an rpc.Client
func NewRegistryRequests(r repo.Repo, cli *rpc.Client) *RegistryRequests {
	if r != nil && cli != nil {
		panic(fmt.Errorf("both repo and client supplied to NewRegistryRequests"))
	}

	return &RegistryRequests{
		repo: actions.Registry{r},
		cli:  cli,
	}
}

// SetQriNode assigns the unexported qri node pointer
func (r *RegistryRequests) SetQriNode(node *p2p.QriNode) {
	r.node = node
}

// PublishParams encapsulates arguments to the publish method
type PublishParams struct {
	Ref repo.DatasetRef
	Pin bool
}

// Publish a dataset to a registry
func (r *RegistryRequests) Publish(p *PublishParams, done *bool) (err error) {
	if r.cli != nil {
		return r.cli.Call("RegistryRequests.Publish", p, done)
	}

	ref := p.Ref

	if p.Pin {
		log.Info("pinning dataset...")
		node := r.node

		if node == nil {
			// if we don't have an online node, create one and connect
			node, err = p2p.NewQriNode(r.repo.Repo, func(c *config.P2P) {
				*c = *Config.P2P
				c.Enabled = true
			})
			if err != nil {
				return err
			}

			if err := node.StartOnlineServices(func(string) {}); err != nil {
				return err
			}
		} else if !node.Online {
			if err := node.StartOnlineServices(func(string) {}); err != nil {
				return err
			}
		}

		var addrs []string
		for _, maddr := range node.EncapsulatedAddresses() {
			addrs = append(addrs, maddr.String())
		}

		if err = r.repo.Pin(ref, addrs); err != nil {
			if err == registry.ErrPinsetNotSupported {
				log.Info("this registry does not support pinning, dataset not pinned.")
			} else {
				return err
			}
		} else {
			log.Info("done")
		}
	}

	return r.repo.Publish(ref)
}

// Unpublish a dataset from a registry
func (r *RegistryRequests) Unpublish(ref *repo.DatasetRef, done *bool) error {
	if r.cli != nil {
		return r.cli.Call("RegistryRequests.Unpublish", ref, done)
	}
	return r.repo.Unpublish(*ref)
}
