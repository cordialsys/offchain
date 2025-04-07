package httpsignature

import (
	"fmt"
	"net/http"
	"strings"
)

type SigBase struct {
	*SigParams
	Method        string
	Path          string
	Query         string
	ContentDigest *ContentDigest
	Headers       http.Header
}

func NewSigBase(params *SigParams, method string, path string, query string, headers http.Header, body []byte) *SigBase {
	content := NewContentDigest(body)

	return &SigBase{
		SigParams:     params,
		Method:        method,
		Path:          path,
		Query:         query,
		ContentDigest: content,
		Headers:       headers,
	}
}

func SigbaseFrom(params *SigParams, method string, path string, query string, headers http.Header, digestHeader *ContentDigest) *SigBase {
	return &SigBase{
		SigParams:     params,
		Method:        method,
		Path:          path,
		Query:         query,
		Headers:       headers,
		ContentDigest: digestHeader,
	}
}

func (s *SigBase) Serialize() (string, error) {
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

	signBase := ""

	for _, component := range s.Components {
		switch component {
		case "@method":
			signBase += fmt.Sprintf(`"@method": %s`, s.Method)
		case "@path":
			signBase += fmt.Sprintf(`"@path": %s`, s.Path)
		case "@query":
			signBase += fmt.Sprintf(`"@query": %s`, s.Query)
		case "content-digest":
			signBase += fmt.Sprintf(`"content-digest": %s`, s.ContentDigest.Base64())
		case "@signature-params":
			signBase += fmt.Sprintf(`"@signature-params": %s`, s.SigParams.Serialize())
		default:
			if strings.HasPrefix(component, "@") {
				return "", fmt.Errorf("unsupported component: %s", component)
			}
			signBase += fmt.Sprintf(`"%s": %s`, component, s.Headers.Get(component))
		}
		// each line has '\n' newline
		signBase += "\n"
	}

	// 	// each line has '\n' newline
	// 	template := fmt.Sprintf(`"@method": %s
	// "@path": %s
	// "@query": %s
	// content-digest: sha-256=:%s:
	// "@signature-params": %s
	// `,
	// 		s.Method,
	// 		path,
	// 		query,
	// 		s.ContentDigest.Base64(),
	// 		s.SigParams.Serialize(),
	// 	)

	return signBase, nil
}
