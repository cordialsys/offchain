package httpsignature

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/cordialsys/offchain/pkg/httpsignature/verifier"
)

func Verify(req *http.Request, verifier verifier.VerifierI) error {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	if req.Header.Get(HeaderContentDigest) == "" {
		return fmt.Errorf("missing content-digest header")
	}
	if req.Header.Get(HeaderSignature) == "" {
		return fmt.Errorf("missing signature header")
	}
	if req.Header.Get(HeaderSignatureInput) == "" {
		return fmt.Errorf("missing signature-input header")
	}

	contentDigest, err := ParseContentDigest(req.Header.Get(HeaderContentDigest))
	if err != nil {
		return err
	}

	signature, err := ParseSignature(req.Header.Get(HeaderSignature))
	if err != nil {
		return err
	}

	signatureInput, err := ParseSigParams(req.Header.Get(HeaderSignatureInput))
	if err != nil {
		return err
	}

	recalculatedContentDigest := NewContentDigest(bodyBytes)
	if !bytes.Equal(recalculatedContentDigest.Digest, contentDigest.Digest) {
		return fmt.Errorf("content-digest mismatch")
	}

	sigBase := SigbaseFrom(signatureInput, req.Method, req.URL.Path, req.URL.RawQuery, contentDigest)
	message := sigBase.Serialize()

	ok := verifier.Verify([]byte(message), signature.Signature)
	if !ok {
		return fmt.Errorf("signature invalid")
	}

	return nil
}
