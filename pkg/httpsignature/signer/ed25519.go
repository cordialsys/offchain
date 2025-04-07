package signer

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

var _ SignerI = &Ed25519Signer{}

type Ed25519Signer struct {
	// privateKey []byte
	Key ed25519.PrivateKey
}

func NewEd25519Signer(secret string) (*Ed25519Signer, error) {
	secret = fmt.Sprintf("%064s", secret)
	seedBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("invalid hex: %v", err)
	}
	if len(seedBytes) != 32 {
		return nil, fmt.Errorf("ed25519 secret not length 32")
	}
	key := ed25519.NewKeyFromSeed(seedBytes)
	return &Ed25519Signer{
		Key: key,
	}, nil
}

func (sk *Ed25519Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(sk.Key, msg), nil
}

func (s *Ed25519Signer) PublicKey() []byte {
	return s.Key.Public().(ed25519.PublicKey)
}
