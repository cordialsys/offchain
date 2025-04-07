package api

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Client struct {
	apiKey     string // Base64 encoded public key
	privateKey ed25519.PrivateKey
	baseURL    string
	httpClient *http.Client
	window     int64 // Time window in milliseconds
}

// NewClient creates a new Backpack API client
func NewClient(apiKey string, privateKey ed25519.PrivateKey) (*Client, error) {
	return &Client{
		apiKey:     apiKey,
		privateKey: privateKey,
		baseURL:    "https://api.backpack.exchange",
		httpClient: &http.Client{},
		window:     5000, // Default window value
	}, nil
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

func (c *Client) SetWindow(window int64) {
	if window > 60000 {
		window = 60000 // Maximum allowed window
	}
	c.window = window
}

// createSigningString creates the string to be signed for authentication
func createSigningString(instruction string, params map[string]string, timestamp int64, window int64) string {
	// Sort keys alphabetically
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string from sorted parameters
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}

	// Create the signing string with instruction, params, timestamp and window
	signingString := fmt.Sprintf("instruction=%s", instruction)
	if len(parts) > 0 {
		signingString += "&" + strings.Join(parts, "&")
	}
	signingString += fmt.Sprintf("&timestamp=%d&window=%d", timestamp, window)

	return signingString
}

// sign creates the signature for authentication
func (c *Client) sign(instruction string, params map[string]string, timestamp int64) string {
	signingString := createSigningString(instruction, params, timestamp, c.window)
	fmt.Println("signingString", signingString)
	signature := ed25519.Sign(c.privateKey, []byte(signingString))
	return base64.StdEncoding.EncodeToString(signature)
}

// Request makes an authenticated HTTP request to the Backpack API
func (c *Client) Request(method, path, instruction string, input interface{}, output interface{}, query url.Values) ([]byte, error) {
	method = strings.ToUpper(method)
	apiUrl := c.baseURL + path

	log := slog.With("method", method, "url", apiUrl, "instruction", instruction)

	// Prepare parameters for signing
	params := make(map[string]string)

	// Handle input body parameters
	var bodyStr string
	if input != nil {
		jsonBody, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyStr = string(jsonBody)

		// Parse JSON body into map for signing
		var bodyMap map[string]interface{}
		if err := json.Unmarshal(jsonBody, &bodyMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal request body for signing: %w", err)
		}

		// Convert all values to strings for the signing map
		for k, v := range bodyMap {
			params[k] = fmt.Sprintf("%v", v)
		}
	}

	// Handle query parameters
	for k, values := range query {
		if len(values) > 0 {
			params[k] = values[0]
		}
	}

	log.Debug("request", "body", bodyStr, "params", params)

	// Generate timestamp
	timestamp := time.Now().UnixMilli()

	// Generate signature
	signature := c.sign(instruction, params, timestamp)

	// Create request
	var reqBody io.Reader
	if bodyStr != "" {
		reqBody = strings.NewReader(bodyStr)
	}

	// Append query to URL if needed
	if len(query) > 0 {
		apiUrl += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, apiUrl, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Window", fmt.Sprintf("%d", c.window))
	req.Header.Set("X-Signature", signature)

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
		var backpackError struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(respBody, &backpackError); err == nil {
			return nil, fmt.Errorf("request failed with code %d: %s", backpackError.Code, backpackError.Message)
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
