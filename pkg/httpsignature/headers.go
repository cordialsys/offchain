package httpsignature

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

const HeaderContentDigest = "content-digest"
const HeaderSignature = "signature"
const HeaderSignatureInput = "signature-input"

type Header interface {
	// returns a tuple (key, value)
	Header() [2]string
}

type ContentDigest struct {
	Algorithm string
	Digest    []byte
}

func NewContentDigest(body []byte) *ContentDigest {
	digest := sha256.Sum256(body)
	return &ContentDigest{
		Algorithm: "sha-256",
		Digest:    digest[:],
	}
}
func (c *ContentDigest) Base64() string {
	return base64.StdEncoding.EncodeToString(c.Digest)
}

func (c *ContentDigest) Header() []string {
	return []string{
		HeaderContentDigest,
		fmt.Sprintf(`sha-256=:%s:`, c.Base64()),
	}
}
func ParseContentDigest(header string) (*ContentDigest, error) {
	eqIndex := strings.Index(header, "=")
	if eqIndex == -1 {
		return nil, fmt.Errorf("invalid content-digest header, expected '='")
	}
	algorithm := header[:eqIndex]
	digestEncoded := header[eqIndex+1:]

	digest, err := base64.StdEncoding.DecodeString(
		strings.Trim(digestEncoded, ":"),
	)
	if err != nil {
		return nil, fmt.Errorf("content-digest invalid base64: %w", err)
	}
	return &ContentDigest{Algorithm: algorithm, Digest: digest}, nil
}

type SigParamKV struct {
	Key   string
	Value string
}

type SigParams struct {
	Name string
	// inside the (...)
	Components []string
	// after the (...);
	Attributes []SigParamKV
}

func NewSigParams(Created int64) *SigParams {
	return &SigParams{
		Name: "iam",
		Components: []string{
			"@method",
			"@path",
			"@query",
			"content-digest",
		},
		Attributes: []SigParamKV{
			{
				Key:   "created",
				Value: fmt.Sprintf("\"%d\"", Created),
			},
		},
	}
}

func (p *SigParams) Serialize() string {
	attributes := []string{}
	for _, attr := range p.Attributes {
		attributes = append(attributes, fmt.Sprintf(`%s=%s`, attr.Key, attr.Value))
	}

	template := fmt.Sprintf(
		`(%s);%s`,
		strings.Join(p.Components, " "),
		strings.Join(attributes, ";"),
	)
	return template
}

func (p *SigParams) Header() []string {
	return []string{
		HeaderSignatureInput,
		fmt.Sprintf(`%s=%s`, p.Name, p.Serialize()),
	}
}
func ParseSigParams(header string) (*SigParams, error) {

	eqIndex := strings.Index(header, "=")
	if eqIndex == -1 {
		return nil, fmt.Errorf("invalid signature-input header, expected '='")
	}
	name := header[:eqIndex]
	value := header[eqIndex+1:]

	paranth0 := strings.Index(value, "(")
	paranth1 := strings.Index(value, ")")
	if paranth0 == -1 || paranth1 == -1 || paranth0 > paranth1 {
		return nil, fmt.Errorf("invalid signature-input header (missing parentheses)")
	}
	components := strings.Split(value[paranth0+1:paranth1], " ")
	attributes := []SigParamKV{}
	kvParts := strings.Split(value[paranth1+1:], ";")
	for _, part := range kvParts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		subParts := strings.Split(part, "=")
		if len(subParts) != 2 {
			return nil, fmt.Errorf("invalid signature-input header (invalid key-value pair %s)", part)
		}
		attributes = append(attributes, SigParamKV{Key: subParts[0], Value: subParts[1]})
	}

	return &SigParams{name, components, attributes}, nil
}

type Signature struct {
	Name      string
	Signature []byte
}

func NewSignature(signature []byte) *Signature {
	return &Signature{
		Name:      "iam",
		Signature: signature,
	}
}

func (sig *Signature) Base64() string {
	return base64.StdEncoding.EncodeToString(sig.Signature)
}

func (sig *Signature) Header() []string {
	return []string{
		HeaderSignature,
		fmt.Sprintf(`%s=:%s:`, sig.Name, sig.Base64()),
	}
}

func ParseSignature(header string) (*Signature, error) {
	eqIndex := strings.Index(header, "=")
	if eqIndex == -1 {
		return nil, fmt.Errorf("invalid signature header, expected '='")
	}
	name := header[:eqIndex]
	value := header[eqIndex+1:]
	signature, err := base64.StdEncoding.DecodeString(
		strings.Trim(value, ":"),
	)
	if err != nil {
		return nil, fmt.Errorf("signature invalid base64: %w", err)
	}
	return &Signature{name, signature}, nil
}
