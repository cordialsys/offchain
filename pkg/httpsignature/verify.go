package httpsignature

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/cordialsys/offchain/pkg/httpsignature/verifier"
	"github.com/gofiber/fiber/v2"
)

type VerifyArgs struct {
	Body           []byte
	Method         string
	Path           string
	Query          string
	Signature      string
	SignatureInput string
	ContentDigest  string
	Headers        http.Header
}

func VerifyRaw(args VerifyArgs, verifier verifier.VerifierI, requiredHeaders ...string) (signatureInput *SigParams, err error) {
	contentDigest, err := ParseContentDigest(args.ContentDigest)
	if err != nil {
		return nil, err
	}

	signature, err := ParseSignature(args.Signature)
	if err != nil {
		return nil, err
	}

	signatureInput, err = ParseSigParams(args.SignatureInput)
	if err != nil {
		return nil, err
	}

	hasPath := false
	hasMethod := false
	hasQuery := false

	for _, component := range signatureInput.Components {
		switch component {
		case "@path":
			hasPath = true
		case "@method":
			hasMethod = true
		case "@query":
			hasQuery = true
		}
	}
	if !hasPath {
		return nil, fmt.Errorf("missing path in signature input")
	}
	if !hasMethod {
		return nil, fmt.Errorf("missing method in signature input")
	}
	if !hasQuery {
		return nil, fmt.Errorf("missing query in signature input")
	}

	for _, header := range requiredHeaders {
		hasHeader := false
		for _, component := range signatureInput.Components {
			if component == header {
				hasHeader = true
				break
			}
		}
		if !hasHeader {
			return nil, fmt.Errorf("missing required header in signature input: %s", header)
		}
	}

	recalculatedContentDigest := NewContentDigest(args.Body)
	if !bytes.Equal(recalculatedContentDigest.Digest, contentDigest.Digest) {
		return nil, fmt.Errorf("content-digest mismatch")
	}

	sigBase := SigbaseFrom(signatureInput, args.Method, args.Path, args.Query, args.Headers, contentDigest)
	message, err := sigBase.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signature base: %w", err)
	}

	ok := verifier.Verify([]byte(message), signature.Signature)
	if !ok {
		return nil, fmt.Errorf("signature invalid")
	}

	return signatureInput, nil
}

func Verify(req *http.Request, verifier verifier.VerifierI, requiredHeaders ...string) (signatureInput *SigParams, err error) {
	// var bodyBytes []byte
	var verifyArgs VerifyArgs
	if req.Body != nil {
		verifyArgs.Body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(verifyArgs.Body))
	}
	if req.Header.Get(HeaderContentDigest) == "" {
		return nil, fmt.Errorf("missing content-digest header")
	}
	if req.Header.Get(HeaderSignature) == "" {
		return nil, fmt.Errorf("missing signature header")
	}
	if req.Header.Get(HeaderSignatureInput) == "" {
		return nil, fmt.Errorf("missing signature-input header")
	}

	verifyArgs.Signature = req.Header.Get(HeaderSignature)
	verifyArgs.SignatureInput = req.Header.Get(HeaderSignatureInput)
	verifyArgs.ContentDigest = req.Header.Get(HeaderContentDigest)
	verifyArgs.Method = req.Method
	verifyArgs.Path = req.URL.Path
	verifyArgs.Query = req.URL.RawQuery
	verifyArgs.Headers = req.Header

	return VerifyRaw(verifyArgs, verifier, requiredHeaders...)
}

func VerifyFiber(c *fiber.Ctx, verifier verifier.VerifierI, requiredHeaders ...string) (signatureInput *SigParams, err error) {
	// var bodyBytes []byte
	var verifyArgs VerifyArgs
	verifyArgs.Body = c.Body()
	if c.Get(HeaderContentDigest) == "" {
		return nil, fmt.Errorf("missing content-digest header")
	}
	if c.Get(HeaderSignature) == "" {
		return nil, fmt.Errorf("missing signature header")
	}
	if c.Get(HeaderSignatureInput) == "" {
		return nil, fmt.Errorf("missing signature-input header")
	}

	verifyArgs.Signature = c.Get(HeaderSignature)
	verifyArgs.SignatureInput = c.Get(HeaderSignatureInput)
	verifyArgs.ContentDigest = c.Get(HeaderContentDigest)
	verifyArgs.Method = c.Method()
	verifyArgs.Path = c.Path()
	verifyArgs.Query = string(c.Request().URI().QueryString())
	verifyArgs.Headers = c.GetReqHeaders()

	return VerifyRaw(verifyArgs, verifier, requiredHeaders...)
}
