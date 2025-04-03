package httpsignature_test

import (
	"encoding/base64"
	"testing"

	"github.com/cordialsys/offchain/pkg/httpsignature"
	"github.com/stretchr/testify/require"
)

func TestContentDigest(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected string
	}{
		{
			name:     "empty body",
			body:     []byte{},
			expected: "47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
		{
			name:     "simple body",
			body:     []byte("test"),
			expected: "n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg=",
		},
		{
			name:     "longer body",
			body:     []byte("This is a longer test message with some content."),
			expected: "6KKSzuCXVyYT0ALX3CS8L0ozLJ1JFTe7B+uBPplNaWU=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			digest := httpsignature.NewContentDigest(tt.body)

			// Test Base64 encoding
			require.Equal(t, tt.expected, digest.Base64())

			// Test Header method
			header := digest.Header()
			require.Equal(t, httpsignature.HeaderContentDigest, header[0])

			expectedHeader := "sha-256=:" + tt.expected + ":"
			require.Equal(t, expectedHeader, header[1])

			// Test parsing
			parsedDigest, err := httpsignature.ParseContentDigest(header[1])
			require.NoError(t, err)
			require.Equal(t, tt.expected, parsedDigest.Base64())

			require.Equal(t, digest.Algorithm, parsedDigest.Algorithm)
			require.Equal(t, "sha-256", digest.Algorithm)
		})
	}
}

func TestParseContentDigestErrors(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		wantErr string
	}{
		{
			name:    "invalid format",
			header:  "invalid-format",
			wantErr: "invalid content-digest header",
		},
		{
			name:    "invalid base64",
			header:  "sha-256=:invalid-base64:",
			wantErr: "content-digest invalid base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := httpsignature.ParseContentDigest(tt.header)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestSigParams(t *testing.T) {
	tests := []struct {
		name       string
		created    int64
		components []string
		attributes []httpsignature.SigParamKV
		expected   string
	}{
		{
			name:     "standard params",
			created:  1234567890,
			expected: "(@method @path @query content-digest);created=\"1234567890\"",
		},
		{
			name:     "zero timestamp",
			created:  0,
			expected: "(@method @path @query content-digest);created=\"0\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := httpsignature.NewSigParams(tt.created)

			// Test serialization
			serialized := params.Serialize()
			require.Equal(t, tt.expected, serialized)

			// Test Header method
			header := params.Header()
			require.Equal(t, httpsignature.HeaderSignatureInput, header[0])

			expectedHeader := "iam=" + tt.expected
			require.Equal(t, expectedHeader, header[1])

			// Test parsing
			parsedParams, err := httpsignature.ParseSigParams(header[1])
			require.NoError(t, err, header[1])
			require.Equal(t, "iam", parsedParams.Name)

			require.Equal(t, params.Components, parsedParams.Components)
			require.Equal(t, params.Attributes, parsedParams.Attributes)
		})
	}
}

func TestParseSigParamsErrors(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		wantErr string
	}{
		{
			name:    "invalid format",
			header:  "invalid-format",
			wantErr: "invalid signature-input header",
		},
		{
			name:    "missing parentheses",
			header:  "iam=no-parentheses;created=\"1234567890\"",
			wantErr: "invalid signature-input header (missing parentheses)",
		},
		{
			name:    "invalid key-value pair",
			header:  "iam=(@method);invalid-kv",
			wantErr: "invalid signature-input header (invalid key-value pair",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := httpsignature.ParseSigParams(tt.header)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestSignature(t *testing.T) {
	tests := []struct {
		name      string
		signature []byte
	}{
		{
			name:      "empty signature",
			signature: []byte{},
		},
		{
			name:      "simple signature",
			signature: []byte("test-signature"),
		},
		{
			name:      "binary signature",
			signature: []byte{0x01, 0x02, 0x03, 0x04, 0xFF, 0xFE, 0xFD, 0xFC},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig := httpsignature.NewSignature(tt.signature)

			// Test Base64 encoding
			expected := base64.StdEncoding.EncodeToString(tt.signature)
			require.Equal(t, expected, sig.Base64())

			// Test Header method
			header := sig.Header()
			require.Equal(t, httpsignature.HeaderSignature, header[0])

			expectedHeader := "iam=:" + expected + ":"
			require.Equal(t, expectedHeader, header[1])

			// Test parsing
			parsedSig, err := httpsignature.ParseSignature(header[1])
			require.NoError(t, err)
			require.Equal(t, "iam", parsedSig.Name)
			require.Equal(t, tt.signature, parsedSig.Signature)
		})
	}
}

func TestParseSignatureErrors(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		wantErr string
	}{
		{
			name:    "invalid format",
			header:  "invalid-format",
			wantErr: "invalid signature header",
		},
		{
			name:    "invalid base64",
			header:  "iam=:invalid-base64:",
			wantErr: "signature invalid base64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := httpsignature.ParseSignature(tt.header)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
