package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"reflect"

	ma "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	peer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	crypto "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"

	"github.com/qri-io/jsonschema"
)

// P2P encapsulates configuration options for qri peer-2-peer communication
type P2P struct {
	// Enabled is a flag for weather this node should connect
	// to the distributed network
	Enabled bool `json:"enabled"`

	// PeerID is this nodes peer identifier
	PeerID string `json:"peerid"`

	PubKey  string `json:"pubkey"`
	PrivKey string `json:"privkey"`

	// Port default port to bind a tcp listener to
	// ignored if Addrs is supplied
	Port int `json:"port"`

	// List of multiaddresses to listen on
	Addrs []ma.Multiaddr `json:"addrs"`

	// QriBootstrapAddrs lists addresses to bootstrap qri node from
	QriBootstrapAddrs []string `json:"qribootstrapaddrs"`

	// HTTPGatewayAddr is an address that qri can use to resolve p2p assets
	// over HTTP, represented as a url. eg: https://ipfs.io
	HTTPGatewayAddr string `json:"httpgatewayaddr"`

	// ProfileReplication determines what to do when this peer sees messages
	// broadcast by it's own profile (from another peer instance). setting
	// ProfileReplication == "full" will cause this peer to automatically pin
	// any data that is verifyably posted by the same peer
	ProfileReplication string `json:"profilereplication"`

	// list of addresses to bootsrap qri peers on
	BootstrapAddrs []string `json:"bootstrapaddrs"`
}

// DefaultP2P generates sensible settings for p2p, generating a new randomized
// private key & peer id
func DefaultP2P() *P2P {
	p := NewP2P()
	_ = p.GeneratePrivateKeyAndPeerID()
	return p
}

// DefaultP2PWithoutKeys generates a p2p struct without keys or peerID
func DefaultP2PWithoutKeys() *P2P {
	return NewP2P()
}

// NewP2P generates a p2p struct with only addresses, no keys or peer id
func NewP2P() *P2P {
	p2p := &P2P{
		Enabled:         true,
		HTTPGatewayAddr: "https://ipfs.io",
		// DefaultBootstrapAddresses follows the pattern of IPFS boostrapping off known "gateways".
		// This boostrapping is specific to finding qri peers, which are IPFS peers that also
		// support the qri protocol.
		// (we also perform standard IPFS boostrapping when IPFS networking is enabled, and it's almost always enabled).
		// These are addresses to public qri nodes hosted by qri, inc.
		// One day it would be super nice to bootstrap from a stored history & only
		// use these for first-round bootstrapping.
		QriBootstrapAddrs: []string{
			"/ip4/130.211.198.23/tcp/4001/ipfs/QmNX9nSos8sRFvqGTwdEme6LQ8R1eJ8EuFgW32F9jjp2Pb", // mojo
			"/ip4/35.193.162.149/tcp/4001/ipfs/QmTZxETL4YCCzB1yFx4GT1te68henVHD1XPQMkHZ1N22mm", // epa
			"/ip4/35.226.92.45/tcp/4001/ipfs/QmP6sbnHXANXgQ7JeCCeCKdJrgpvUd8s75YNfzdkHf6Mpi",   // 538
			"/ip4/35.192.140.245/tcp/4001/ipfs/QmUUVNiTz2K9zQSH9PxerKWXmN1p3DBo3oJXurvYziFzqh", // EDGI
		},
		ProfileReplication: "full",
	}
	return p2p
}

// GeneratePrivateKeyAndPeerID generates a new random private key and peer id
func (cfg *P2P) GeneratePrivateKeyAndPeerID() error {
	r := rand.Reader
	// Generate a key pair for this host
	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return err
	}
	pdata, err := priv.Bytes()
	if err != nil {
		return err
	}
	cfg.PrivKey = base64.StdEncoding.EncodeToString(pdata)
	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return err
	}
	cfg.PeerID = pid.Pretty()
	return nil
}

// DecodePrivateKey generates a PrivKey instance from base64-encoded config file bytes
func (cfg *P2P) DecodePrivateKey() (crypto.PrivKey, error) {
	if cfg.PrivKey == "" {
		return nil, fmt.Errorf("missing private key")
	}

	data, err := base64.StdEncoding.DecodeString(cfg.PrivKey)
	if err != nil {
		return nil, fmt.Errorf("decoding private key: %s", err.Error())
	}

	return crypto.UnmarshalPrivateKey(data)
}

// DecodePeerID takes P2P.ID (a string), and decodes it into a peer.ID
func (cfg *P2P) DecodePeerID() (peer.ID, error) {
	if string(cfg.PeerID) == "" {
		return peer.ID(""), fmt.Errorf("empty string for peer ID")
	}
	return peer.IDB58Decode(cfg.PeerID)
}

// Validate validates all fields of p2p returning all errors found.
func (cfg P2P) Validate() error {
	schema := jsonschema.Must(`{
    "$schema": "http://json-schema.org/draft-06/schema#",
    "title": "P2P",
    "description": "Config for the p2p",
    "type": "object",
    "required": ["enabled", "peerid", "pubkey", "privkey", "port", "addrs", "httpgatewayaddr", "qribootstrapaddrs", "profilereplication", "bootstrapaddrs"],
    "properties": {
      "enabled": {
        "description": "When true, peer to peer communication is allowed",
        "type": "boolean"
      },
      "peerid": {
        "description": "The peerid is this nodes peer identifier",
        "type": "string"
      },
      "pubkey": {
        "description": "",
        "type": "string"
      },
      "privkey": {
        "description": "",
        "type": "string"
      },
      "port": {
        "description": "Port to bind a tcp lister to. Field is ignored if addrs is supplied",
        "type": "integer"
      },
      "addrs": {
        "description": "List of multiaddresses to listen on",
        "anyOf": [
          {"type": "array"},
          {"type": "null"}
        ],
        "items": {
          "type": "string"
        }
      },
      "httpgatewayaddr": {
        "description" : "address that qri can use to resolve p2p assets over HTTP",
        "type" : "string"
      },
      "qribootstrapaddrs": {
        "description": "List of addresses to bootstrap the qri node from",
        "type": "array",
        "items": {
          "type": "string"
        }
      },
      "profilereplication": {
        "description": "Determings what to do when this peer sees messages broadcast by it's own profile (from another peer instance). Setting profilereplication to 'full' will cause this peer to automatically pin any data that is verifiably posted by the same peer",
        "type": "string",
        "enum": [
          "full"
        ]
      },
      "bootstrapaddrs": {
        "description": "List of addresses to bootstrap qri peers on",
        "anyOf": [
          {"type": "array"},
          {"type": "null"}
        ],
        "items": {
          "type": "string"
        }
      }
    }
  }`)
	return validate(schema, &cfg)
}

// Copy returns a deep copy of a p2p struct
func (cfg *P2P) Copy() *P2P {
	res := &P2P{
		Enabled:            cfg.Enabled,
		PeerID:             cfg.PeerID,
		PubKey:             cfg.PubKey,
		PrivKey:            cfg.PrivKey,
		Port:               cfg.Port,
		ProfileReplication: cfg.ProfileReplication,
		HTTPGatewayAddr:    cfg.HTTPGatewayAddr,
	}

	if cfg.QriBootstrapAddrs != nil {
		res.QriBootstrapAddrs = make([]string, len(cfg.QriBootstrapAddrs))
		reflect.Copy(reflect.ValueOf(res.QriBootstrapAddrs), reflect.ValueOf(cfg.QriBootstrapAddrs))
	}

	if cfg.BootstrapAddrs != nil {
		res.BootstrapAddrs = make([]string, len(cfg.BootstrapAddrs))
		reflect.Copy(reflect.ValueOf(res.BootstrapAddrs), reflect.ValueOf(cfg.BootstrapAddrs))
	}

	return res
}
