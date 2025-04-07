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
}

func VerifyRaw(args VerifyArgs, verifier verifier.VerifierI) error {
	contentDigest, err := ParseContentDigest(args.ContentDigest)
	if err != nil {
		return err
	}

	signature, err := ParseSignature(args.Signature)
	if err != nil {
		return err
	}

	signatureInput, err := ParseSigParams(args.SignatureInput)
	if err != nil {
		return err
	}

	recalculatedContentDigest := NewContentDigest(args.Body)
	if !bytes.Equal(recalculatedContentDigest.Digest, contentDigest.Digest) {
		return fmt.Errorf("content-digest mismatch")
	}

	sigBase := SigbaseFrom(signatureInput, args.Method, args.Path, args.Query, contentDigest)
	message := sigBase.Serialize()

	ok := verifier.Verify([]byte(message), signature.Signature)
	if !ok {
		return fmt.Errorf("signature invalid")
	}
	return nil
}

func Verify(req *http.Request, verifier verifier.VerifierI) error {
	// var bodyBytes []byte
	var verifyArgs VerifyArgs
	if req.Body != nil {
		verifyArgs.Body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(verifyArgs.Body))
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

	verifyArgs.Signature = req.Header.Get(HeaderSignature)
	verifyArgs.SignatureInput = req.Header.Get(HeaderSignatureInput)
	verifyArgs.ContentDigest = req.Header.Get(HeaderContentDigest)
	verifyArgs.Method = req.Method
	verifyArgs.Path = req.URL.Path
	verifyArgs.Query = req.URL.RawQuery

	return VerifyRaw(verifyArgs, verifier)
}

func VerifyFiber(c *fiber.Ctx, verifier verifier.VerifierI) error {
	// var bodyBytes []byte
	var verifyArgs VerifyArgs
	verifyArgs.Body = c.Body()
	if c.Get(HeaderContentDigest) == "" {
		return fmt.Errorf("missing content-digest header")
	}
	if c.Get(HeaderSignature) == "" {
		return fmt.Errorf("missing signature header")
	}
	if c.Get(HeaderSignatureInput) == "" {
		return fmt.Errorf("missing signature-input header")
	}

	verifyArgs.Signature = c.Get(HeaderSignature)
	verifyArgs.SignatureInput = c.Get(HeaderSignatureInput)
	verifyArgs.ContentDigest = c.Get(HeaderContentDigest)
	verifyArgs.Method = c.Method()
	verifyArgs.Path = c.Path()
	verifyArgs.Query = string(c.Request().URI().QueryString())

	return VerifyRaw(verifyArgs, verifier)
}
