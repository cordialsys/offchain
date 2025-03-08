package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL = "https://api.binance.us"
)

type Client struct {
	apiKey     string
	secretKey  string
	httpClient *http.Client
	baseURL    string
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		secretKey:  secretKey,
		httpClient: &http.Client{},
		baseURL:    BaseURL,
	}
}

// signRequest creates HMAC SHA256 signature for the request
func (c *Client) signRequest(queryString string, body []byte) string {
	// Create HMAC SHA256 signature
	h := hmac.New(sha256.New, []byte(c.secretKey))
	fmt.Println("queryString", queryString)
	fmt.Println("body", string(body))
	h.Write([]byte(queryString))
	h.Write(body)
	signature := hex.EncodeToString(h.Sum(nil))
	return signature
}

// createQueryString creates a sorted query string with timestamp
func (c *Client) createQueryString(query url.Values) string {
	if query == nil {
		query = url.Values{}
	}

	// Add timestamp if not present
	if query.Get("timestamp") == "" {
		query.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	// Sort query parameters
	params := make([]string, 0, len(query))
	for key := range query {
		params = append(params, key)
	}
	sort.Strings(params)

	// Build sorted query string
	var builder strings.Builder
	for i, key := range params {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(query.Get(key))
	}
	return builder.String()
}

func (c *Client) Request(method, path string, body interface{}, output interface{}, query url.Values) ([]byte, error) {
	endpoint := c.baseURL + path
	var reqBody []byte
	var err error

	// Convert body to JSON if present
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create query string and signature
	queryString := c.createQueryString(query)
	signature := c.signRequest(queryString, reqBody)

	// Add signature to query parameters
	if query == nil {
		query = url.Values{}
	}
	query.Set("signature", signature)

	// Create full URL with query parameters
	fullURL := endpoint
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	// Create request
	req, err := http.NewRequest(method, fullURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-MBX-APIKEY", c.apiKey)
	if len(reqBody) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	slog.Debug("sending request",
		"method", method,
		"url", fullURL,
		"body", string(reqBody),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Debug("received response",
		"status", resp.StatusCode,
		"body", string(respBody),
	)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed %d: %s", resp.StatusCode, string(respBody))
	}

	// Unmarshal response if output interface provided
	if output != nil {
		err = json.Unmarshal(respBody, output)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return respBody, nil
}
