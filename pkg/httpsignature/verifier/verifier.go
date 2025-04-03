package verifier

type VerifierI interface {
	Verify(message []byte, signature []byte) bool
}
