package signer

type SignerI interface {
	Sign(data []byte) ([]byte, error)
	PublicKey() []byte
}
