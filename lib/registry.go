package lib

import (
	"fmt"
	"net/rpc"

	"github.com/qri-io/qri/actions"
	"github.com/qri-io/qri/p2p"
	"github.com/qri-io/qri/repo"
)

// RegistryRequests defines business logic for working with registries
type RegistryRequests struct {
	node *p2p.QriNode
	cli  *rpc.Client
}

// CoreRequestsName implements the Requests interface
func (RegistryRequests) CoreRequestsName() string { return "registry" }

// NewRegistryRequests creates a RegistryRequests pointer from either a repo
// or an rpc.Client
func NewRegistryRequests(node *p2p.QriNode, cli *rpc.Client) *RegistryRequests {
	if node != nil && cli != nil {
		panic(fmt.Errorf("both repo and client supplied to NewRegistryRequests"))
	}

	return &RegistryRequests{
		node: node,
		cli:  cli,
	}
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

	// if p.Pin {
	// 	log.Info("pinning dataset...")
	// 	node := r.node

	// 	if !node.Online {
	// 		if err = node.Connect(); err != nil {
	// 			return err
	// 		}
	// 		if err = node.StartOnlineServices(func(string) {}); err != nil {
	// 			return err
	// 		}
	// 	}

	// 	var addrs []string
	// 	for _, maddr := range node.EncapsulatedAddresses() {
	// 		addrs = append(addrs, maddr.String())
	// 	}

	// 	if err = r.node.repo.Pin(ref, addrs); err != nil {
	// 		if err == registry.ErrPinsetNotSupported {
	// 			log.Info("this registry does not support pinning, dataset not pinned.")
	// 		} else {
	// 			return err
	// 		}
	// 	} else {
	// 		log.Info("done")
	// 	}
	// }

	return actions.Publish(r.node, ref)
}

// Unpublish a dataset from a registry
func (r *RegistryRequests) Unpublish(ref *repo.DatasetRef, done *bool) error {
	if r.cli != nil {
		return r.cli.Call("RegistryRequests.Unpublish", ref, done)
	}
	return actions.Unpublish(r.node, *ref)
}

// Status checks if a dataset has been published to a registry
func (r *RegistryRequests) Status(ref *repo.DatasetRef, done *bool) error {
	if r.cli != nil {
		return r.cli.Call("RegistryRequests.Status", ref, done)
	}
	return actions.Status(r.node, *ref)
}
