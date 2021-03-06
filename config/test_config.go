package config

import (
	"github.com/qri-io/qri/config/test"
)

// DefaultConfigForTesting constructs a config with no keys, only used for testing.
func DefaultConfigForTesting() *Config {
	info := test.GetTestPeerInfo(0)
	cfg := DefaultConfigWithoutKeys()
	cfg.P2P.PrivKey = info.EncodedPrivKey
	cfg.P2P.PeerID = info.EncodedPeerID
	cfg.Profile.PrivKey = info.EncodedPrivKey
	cfg.Profile.ID = info.EncodedPeerID
	return cfg
}

// DefaultProfileForTesting constructs a profile with no keys, only used for testing.
func DefaultProfileForTesting() *ProfilePod {
	info := test.GetTestPeerInfo(0)
	pro := DefaultProfileWithoutKeys()
	pro.PrivKey = info.EncodedPrivKey
	pro.ID = info.EncodedPeerID
	return pro
}

// DefaultP2PForTesting constructs a p2p with no keys, only used for testing.
func DefaultP2PForTesting() *P2P {
	info := test.GetTestPeerInfo(0)
	p := DefaultP2PWithoutKeys()
	p.PrivKey = info.EncodedPrivKey
	p.PeerID = info.EncodedPeerID
	return p
}
