package httpsignature_test

import (
	"bytes"
	"crypto/ed25519"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/cordialsys/offchain/pkg/httpsignature"
	"github.com/cordialsys/offchain/pkg/httpsignature/signer"
	"github.com/cordialsys/offchain/pkg/httpsignature/verifier"
	"github.com/stretchr/testify/require"
)

type mockVerifier struct {
	shouldVerify bool
}

func (m *mockVerifier) Verify(message []byte, signature []byte) bool {
	return m.shouldVerify
}

func TestVerify(t *testing.T) {
	// Create a fixed time for consistent tests
	originalNow := httpsignature.Now
	httpsignature.Now = func() time.Time {
		return time.Unix(1234567890, 0)
	}
	defer func() {
		httpsignature.Now = originalNow
	}()

	// Create a test key pair
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	// Create a signer with the private key
	testSigner := &signer.Ed25519Signer{
		Key: privKey,
	}

	// Create a verifier with the public key
	testVerifier, err := verifier.NewEd25519Verifier(pubKey)
	require.NoError(t, err)

	tests := []struct {
		name     string
		setupReq func() *http.Request
		verifier verifier.VerifierI
		errorMsg string
	}{
		{
			name: "valid signature",
			setupReq: func() *http.Request {
				body := []byte("test body")
				req, _ := http.NewRequest("GET", "https://example.com/path?query=value", bytes.NewReader(body))
				err := httpsignature.Sign(req, testSigner)
				require.NoError(t, err)
				return req
			},
			verifier: testVerifier,
		},
		{
			name: "missing content-digest header",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/path", nil)
				// Don't sign the request, so headers will be missing
				return req
			},
			verifier: testVerifier,
			errorMsg: "missing content-digest header",
		},
		{
			name: "missing signature header",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/path", nil)
				// Add only content-digest header
				digest := httpsignature.NewContentDigest([]byte("test"))
				h := digest.Header()
				req.Header.Set(h[0], h[1])
				return req
			},
			verifier: testVerifier,
			errorMsg: "missing signature header",
		},
		{
			name: "missing signature-input header",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/path", nil)
				// Add content-digest and signature headers but not signature-input
				digest := httpsignature.NewContentDigest([]byte("test"))
				h1 := digest.Header()
				req.Header.Set(h1[0], h1[1])

				sig := httpsignature.NewSignature([]byte("fake-signature"))
				h2 := sig.Header()
				req.Header.Set(h2[0], h2[1])
				return req
			},
			verifier: testVerifier,
			errorMsg: "missing signature-input header",
		},
		{
			name: "invalid content-digest format",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "https://example.com/path", nil)
				req.Header.Set(httpsignature.HeaderContentDigest, "invalid-format")
				req.Header.Set(httpsignature.HeaderSignature, "iam=:c2lnbmF0dXJl:")
				req.Header.Set(httpsignature.HeaderSignatureInput, "iam=(@method @path @query content-digest);created=\"1234567890\"")
				return req
			},
			verifier: testVerifier,
			errorMsg: "invalid content-digest header",
		},
		{
			name: "content-digest mismatch",
			setupReq: func() *http.Request {
				body := []byte("test body")
				req, _ := http.NewRequest("GET", "https://example.com/path", bytes.NewReader(body))

				// Sign the request
				err := httpsignature.Sign(req, testSigner)
				require.NoError(t, err)

				// Replace the body with different content
				req.Body = io.NopCloser(bytes.NewReader([]byte("different body")))

				return req
			},
			verifier: testVerifier,
			errorMsg: "content-digest mismatch",
		},
		{
			name: "invalid signature",
			setupReq: func() *http.Request {
				body := []byte("test body")
				req, _ := http.NewRequest("GET", "https://example.com/path", bytes.NewReader(body))

				// Sign the request
				err := httpsignature.Sign(req, testSigner)
				require.NoError(t, err)

				return req
			},
			verifier: &mockVerifier{shouldVerify: false},
			errorMsg: "signature invalid",
		},
		{
			name: "with normalized path and query",
			setupReq: func() *http.Request {
				body := []byte("test body")
				req, _ := http.NewRequest("GET", "https://example.com/path", bytes.NewReader(body))

				// Manually set URL without leading slashes to test normalization
				req.URL = &url.URL{
					Path:     "path",
					RawQuery: "query=value",
				}

				err := httpsignature.Sign(req, testSigner)
				require.NoError(t, err)

				return req
			},
			verifier: testVerifier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()
			err := httpsignature.Verify(req, tt.verifier)

			if tt.errorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
