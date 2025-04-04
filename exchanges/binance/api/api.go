package api

import (
	"crypto/hmac"
	"crypto/sha256"
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
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Binance API client
func NewClient(apiKey, secretKey string) (*Client, error) {
	return &Client{
		apiKey:     apiKey,
		secretKey:  secretKey,
		baseURL:    "https://api.binance.com",
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// sign creates the signature for authentication
func (c *Client) sign(params url.Values) string {
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(params.Encode()))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// Request makes an authenticated HTTP request to the Binance API
func (c *Client) Request(method, path string, input interface{}, output interface{}, query url.Values) ([]byte, error) {
	method = strings.ToUpper(method)
	apiUrl := c.baseURL + path

	log := slog.With("method", method, "url", apiUrl, "query", query)

	// Initialize query parameters if nil
	if query == nil {
		query = url.Values{}
	}

	// Add timestamp for signed endpoints
	query.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// Sign the request
	signature := c.sign(query)
	query.Set("signature", signature)

	// Append query to path
	if len(query) > 0 {
		apiUrl += "?" + query.Encode()
	}

	var bodyStr string
	if input != nil {
		jsonBody, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyStr = string(jsonBody)
	}
	log.Debug("request", "body", bodyStr)

	req, err := http.NewRequest(method, apiUrl, strings.NewReader(bodyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-MBX-APIKEY", c.apiKey)
	fmt.Println("request", req.Header)

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

	if resp.StatusCode != http.StatusOK {
		var binanceError struct {
			Code    int    `json:"code"`
			Message string `json:"msg"`
		}
		if err := json.Unmarshal(respBody, &binanceError); err == nil {
			return nil, fmt.Errorf("request failed with code %d: %s", binanceError.Code, binanceError.Message)
		}
		return nil, fmt.Errorf("request failed %d: %s", resp.StatusCode, string(respBody))
	}

	if output != nil {
		err = json.Unmarshal(respBody, output)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return respBody, nil
}
