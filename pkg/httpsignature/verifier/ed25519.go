package verifier

import (
	"crypto/ed25519"
	"fmt"
)

type Ed25519Verifier struct {
	pubKey []byte
}

func NewEd25519Verifier(pubKey []byte) (*Ed25519Verifier, error) {
	if len(pubKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key length")
	}
	return &Ed25519Verifier{pubKey: pubKey}, nil
}

func (v *Ed25519Verifier) Verify(message []byte, signature []byte) bool {
	return ed25519.Verify(v.pubKey, message, signature)
}
