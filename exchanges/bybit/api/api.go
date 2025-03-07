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
	recvWindow time.Duration
}

// NewClient creates a new OKX API client
func NewClient(apiKey, secretKey string) (*Client, error) {
	return &Client{
		apiKey:     apiKey,
		secretKey:  secretKey,
		baseURL:    "https://api.bybit.com",
		httpClient: &http.Client{},
		recvWindow: time.Second * 5,
	}, nil
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// sign creates the signature for authentication
func (c *Client) sign(timestamp int64, method, queryString, body string) string {
	var payload string
	// Docs / examples say to use miliseconds but it's actually nanoseconds. _confused_
	recvWindowMillis := c.recvWindow.Nanoseconds()
	if method == "GET" {
		payload = fmt.Sprintf("%d%s%d%s", timestamp, c.apiKey, recvWindowMillis, queryString)
	} else {
		payload = fmt.Sprintf("%d%s%d%s", timestamp, c.apiKey, recvWindowMillis, body)
	}

	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(payload))
	return fmt.Sprintf("%x", mac.Sum(nil))
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

	timestamp := time.Now().UnixMilli()
	queryStr := query.Encode()
	signature := c.sign(timestamp, method, queryStr, bodyStr)

	req, err := http.NewRequest(method, apiUrl, strings.NewReader(bodyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for Bybit
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-BAPI-API-KEY", c.apiKey)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-BAPI-RECV-WINDOW", fmt.Sprintf("%d", c.recvWindow))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w (status = %d)", err, resp.StatusCode)
	}

	log.Debug("response", "status", resp.StatusCode, "body", string(respBody))

	if output != nil {
		err = json.Unmarshal(respBody, output)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w (status = %d)", err, resp.StatusCode)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	responseWrapper := Response[json.RawMessage]{}
	err = json.Unmarshal(respBody, &responseWrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w (status = %d)", err, resp.StatusCode)
	}
	if responseWrapper.RetCode != 0 {
		return nil, fmt.Errorf(
			"request failed with application status %s: %s",
			responseWrapper.RetCode.String(),
			responseWrapper.RetMsg,
		)
	}

	return respBody, nil
}
