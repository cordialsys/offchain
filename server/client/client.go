package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cordialsys/offchain/pkg/httpsignature"
	"github.com/cordialsys/offchain/pkg/httpsignature/signer"
	"github.com/cordialsys/offchain/server/client/api"
)

// Client represents a client for the offchain server API
type Client struct {
	baseURL    string
	httpClient *http.Client

	// Authorization options
	bearerToken string
	signer      signer.SignerI
	subAccount  string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBearerToken configures the client to use bearer token authentication
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.bearerToken = token
	}
}

// WithSigner configures the client to use HTTP signature authentication
func WithSigner(signer signer.SignerI) ClientOption {
	return func(c *Client) {
		c.signer = signer
	}
}

func WithSubAccount(subAccount string) ClientOption {
	return func(c *Client) {
		c.subAccount = subAccount
	}
}

// WithTimeout configures the client's HTTP timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new API client
func NewClient(baseURL string, options ...ClientOption) *Client {
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// APIError represents an error returned by the API
type APIError struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Status, e.Message)
}

// doRequest performs an HTTP request with the appropriate authentication and handles response parsing
func (c *Client) doRequest(method, path string, queryParams url.Values, body interface{}, result interface{}) error {
	// Construct the full URL
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	reqURL.Path = path

	query := url.Values{}
	if queryParams != nil {
		query = queryParams
	}
	if c.subAccount != "" {
		query.Set("sub-account", c.subAccount)
	}
	reqURL.RawQuery = query.Encode()

	// Prepare the request body
	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create the request
	req, err := http.NewRequest(method, reqURL.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply authentication
	if method == http.MethodGet && c.bearerToken != "" {
		// For GET requests, prefer bearer token if available
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else if c.signer != nil {
		// For non-GET requests or if bearer token is not available, use HTTP signature
		err = httpsignature.Sign(req, c.signer)
		if err != nil {
			return fmt.Errorf("failed to sign request: %w", err)
		}
	} else if c.bearerToken != "" {
		// Fall back to bearer token for GET requests if HTTP signature is not available
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else {
		return fmt.Errorf("no authentication method provided")
	}

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		// Try to parse as structured error
		var apiError APIError
		bodyBytes, _ := io.ReadAll(resp.Body)

		if err := json.Unmarshal(bodyBytes, &apiError); err == nil && apiError.Message != "" {
			return &apiError
		}

		// Fall back to generic error if parsing fails
		return fmt.Errorf("API error: %s, status code: %d", string(bodyBytes), resp.StatusCode)
	}

	// Parse the response if a result container was provided
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// GetBalances retrieves balances for a specified account
func (c *Client) GetBalances(exchange string, accountType string) ([]*api.Balance, error) {
	queryParams := url.Values{}
	if accountType != "" {
		queryParams.Set("type", accountType)
	}

	var balances []*api.Balance
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/exchanges/%s/balances", exchange), queryParams, nil, &balances)
	return balances, err
}

// GetAssets retrieves the list of assets for a specified exchange
func (c *Client) GetAssets(exchange string) ([]*api.Asset, error) {
	var assets []*api.Asset
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/exchanges/%s/assets", exchange), nil, nil, &assets)
	return assets, err
}

// GetDepositAddress retrieves a deposit address for a specified symbol and network
func (c *Client) GetDepositAddress(exchange, symbol, network string) (string, error) {
	queryParams := url.Values{}
	queryParams.Set("symbol", symbol)
	queryParams.Set("network", network)

	var address string
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/exchanges/%s/deposit-address", exchange), queryParams, nil, &address)
	return address, err
}

// GetAccountTypes retrieves the list of valid account types for an exchange
func (c *Client) GetAccountTypes(exchange string) ([]*api.AccountType, error) {
	var accountTypes []*api.AccountType
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/exchanges/%s/account-types", exchange), nil, nil, &accountTypes)
	return accountTypes, err
}

// ListSubaccounts retrieves the list of configured subaccounts on an exchange
func (c *Client) ListSubaccounts(exchange string) ([]*api.SubAccountHeader, error) {
	var subaccounts []*api.SubAccountHeader
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/exchanges/%s/subaccounts", exchange), nil, nil, &subaccounts)
	return subaccounts, err
}

// ListWithdrawalHistory retrieves the withdrawal history for an exchange account
func (c *Client) ListWithdrawalHistory(exchange string, limit int, pageToken string) ([]*api.HistoricalWithdrawal, error) {
	queryParams := url.Values{}
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if pageToken != "" {
		queryParams.Set("page_token", pageToken)
	}

	var withdrawals []*api.HistoricalWithdrawal
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/exchanges/%s/withdrawal-history", exchange), queryParams, nil, &withdrawals)
	return withdrawals, err
}

// CreateAccountTransfer performs a transfer between accounts on an exchange
func (c *Client) CreateAccountTransfer(exchange string, transfer *api.Transfer) (*api.TransferResponse, error) {
	// HTTP signature is required for this endpoint
	if c.signer == nil {
		return nil, fmt.Errorf("HTTP signature is required for account transfers")
	}

	var transferResp api.TransferResponse
	err := c.doRequest(http.MethodPost, fmt.Sprintf("/v1/exchanges/%s/account-transfer", exchange), nil, transfer, &transferResp)
	if err != nil {
		return nil, err
	}
	return &transferResp, nil
}

// CreateWithdrawal initiates a withdrawal from an exchange
func (c *Client) CreateWithdrawal(exchange string, withdrawal *api.Withdrawal) (*api.WithdrawalResponse, error) {
	// HTTP signature is required for this endpoint
	if c.signer == nil {
		return nil, fmt.Errorf("HTTP signature is required for withdrawals")
	}

	var withdrawalResp api.WithdrawalResponse
	err := c.doRequest(http.MethodPost, fmt.Sprintf("/v1/exchanges/%s/withdrawal", exchange), nil, withdrawal, &withdrawalResp)
	if err != nil {
		return nil, err
	}
	return &withdrawalResp, nil
}
