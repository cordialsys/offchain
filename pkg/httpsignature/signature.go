package httpsignature

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cordialsys/offchain/pkg/httpsignature/signer"
)

var Now = time.Now

func Sign(req *http.Request, signer signer.SignerI) error {
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	created := Now().Unix()
	params := NewSigParams(created)
	signatureBase := NewSigBase(params, req.Method, req.URL.Path, req.URL.RawQuery, bodyBytes)

	sigBaseBz := signatureBase.Serialize()
	rawSig, err := signer.Sign([]byte(sigBaseBz))
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}
	signature := NewSignature(rawSig)

	h1 := signatureBase.SigParams.Header()
	h2 := signatureBase.ContentDigest.Header()
	h3 := signature.Header()
	req.Header.Set(h1[0], h1[1])
	req.Header.Set(h2[0], h2[1])
	req.Header.Set(h3[0], h3[1])
	return nil
}
