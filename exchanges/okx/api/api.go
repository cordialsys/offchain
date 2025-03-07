package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	apiKey     string
	secretKey  string
	passphrase string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new OKX API client
func NewClient(apiKey, secretKey, passphrase string) (*Client, error) {
	return &Client{
		apiKey:     apiKey,
		secretKey:  secretKey,
		passphrase: passphrase,
		baseURL:    "https://www.okx.com",
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// sign creates the signature for authentication
func (c *Client) sign(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Request makes an authenticated HTTP request to the OKX API
func (c *Client) Request(method, path string, input interface{}, output interface{}, query url.Values) ([]byte, error) {
	method = strings.ToUpper(method)
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	apiUrl := c.baseURL + path

	log := slog.With("method", method, "url", apiUrl)

	var bodyStr string
	if input != nil {
		jsonBody, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyStr = string(jsonBody)
	}
	log.Debug("request", "body", bodyStr)

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	signature := c.sign(timestamp, method, path, bodyStr)

	req, err := http.NewRequest(method, apiUrl, strings.NewReader(bodyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OK-ACCESS-KEY", c.apiKey)
	req.Header.Set("OK-ACCESS-SIGN", signature)
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", c.passphrase)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Debug("response", "status", resp.StatusCode, "body", string(respBody))

	if output != nil {
		err = json.Unmarshal(respBody, output)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	responseWrapper := Response[json.RawMessage]{}
	err = json.Unmarshal(respBody, &responseWrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w (status = %d)", err, resp.StatusCode)
	}
	if responseWrapper.Code != "0" {
		return nil, fmt.Errorf(
			"request failed with application status %s: %s",
			responseWrapper.Code,
			responseWrapper.Msg,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
