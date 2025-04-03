package httpsignature

import (
	"fmt"
	"strings"
)

type SigBase struct {
	*SigParams
	Method        string
	Path          string
	Query         string
	ContentDigest *ContentDigest
}

func NewSigBase(params *SigParams, method string, path string, query string, body []byte) *SigBase {
	content := NewContentDigest(body)

	return &SigBase{
		SigParams:     params,
		Method:        method,
		Path:          path,
		Query:         query,
		ContentDigest: content,
	}
}

func SigbaseFrom(params *SigParams, method string, path string, query string, digestHeader *ContentDigest) *SigBase {
	return &SigBase{
		SigParams:     params,
		Method:        method,
		Path:          path,
		Query:         query,
		ContentDigest: digestHeader,
	}
}

func (s *SigBase) Serialize() string {
	query := s.Query
	if !strings.HasPrefix(query, "?") {
		// should have leading ?
		query = "?" + query
	}
	path := s.Path
	if !strings.HasPrefix(path, "/") {
		// should have leading /
		path = "/" + path
	}
	// each line has '\n' newline
	template := fmt.Sprintf(`"@method": %s
"@path": %s
"@query": %s
content-digest: sha-256=:%s:
"@signature-params": %s
`,
		s.Method,
		path,
		query,
		s.ContentDigest.Base64(),
		s.SigParams.Serialize(),
	)
	// fmt.Printf("%s\n", template)

	return template
}
