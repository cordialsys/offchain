package backpack

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strconv"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/backpack/api"
)

type Client struct {
	api *api.Client
}

var _ client.Client = &Client{}

func NewClient(config *oc.ExchangeClientConfig, account *oc.Account) (*Client, error) {
	apiKey, err := account.ApiKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load api key: %w", err)
	}

	secretKey, err := account.SecretKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load secret key: %w", err)
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Decode the private key from base64
	privateKeyBytes, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	if len(privateKeyBytes) != ed25519.SeedSize {
		return nil, fmt.Errorf("invalid secret key size, expected %d bytes, got %d", ed25519.SeedSize, len(privateKeyBytes))
	}

	// Backpack uses ED25519 keys
	if len(privateKeyBytes) != ed25519.SeedSize {
		return nil, fmt.Errorf("invalid private key size, expected %d bytes, got %d", ed25519.SeedSize, len(privateKeyBytes))
	}
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size, expected %d bytes, got %d", ed25519.PublicKeySize, len(publicKeyBytes))
	}

	keyPairFromSeed := ed25519.NewKeyFromSeed(privateKeyBytes)
	pubkeyFromSeed := keyPairFromSeed.Public()
	if !bytes.Equal(pubkeyFromSeed.(ed25519.PublicKey), publicKeyBytes) {
		return nil, fmt.Errorf("the api-secret does not correspond to the api-key")
	}

	apiClient, err := api.NewClient(apiKey, keyPairFromSeed)
	if err != nil {
		return nil, fmt.Errorf("failed to create api client: %w", err)
	}

	return &Client{api: apiClient}, nil
}

func (c *Client) ListAssets() ([]*oc.Asset, error) {
	response, err := c.api.GetAssets()
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}

	assets := make([]*oc.Asset, 0, len(response))
	for _, asset := range response {
		// If there are no tokens, add the asset with just the symbol
		if len(asset.Tokens) == 0 {
			assets = append(assets, &oc.Asset{
				SymbolId: oc.SymbolId(asset.Symbol),
			})
		} else {
			for _, token := range asset.Tokens {
				assets = append(assets, &oc.Asset{
					SymbolId:        oc.SymbolId(asset.Symbol),
					NetworkId:       token.Blockchain,
					ContractAddress: token.ContractAddress,
				})
			}
		}
	}

	return assets, nil
}

func (c *Client) ListBalances(args client.GetBalanceArgs) ([]*client.BalanceDetail, error) {
	response, err := c.api.GetBalances()
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	balances := make([]*client.BalanceDetail, 0, len(response))
	for symbol, balance := range response {
		unavailable := balance.Locked.Decimal().Add(balance.Staked.Decimal())
		balances = append(balances, &client.BalanceDetail{
			SymbolId:    oc.SymbolId(symbol),
			Available:   balance.Available,
			Unavailable: oc.Amount(unavailable),
		})
	}

	return balances, nil
}

// Need some slight additions to Backpack exchange API to support this
func (c *Client) CreateAccountTransfer(args client.AccountTransferArgs) (*client.TransferStatus, error) {
	// // This currently fails since it's not a supported API endpoint
	// mapping, err := c.api.GetSubaccounts()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get subaccounts: %w", err)
	// }

	// debug.PrintJson(mapping)

	return nil, fmt.Errorf("transfers not currently implemented for Backpack exchange")
}

func (c *Client) CreateWithdrawal(args client.WithdrawalArgs) (*client.WithdrawalResponse, error) {
	request := api.WithdrawalRequest{
		Address:    args.GetAddress(),
		Blockchain: args.GetNetwork(),
		Quantity:   args.GetAmount(),
		Symbol:     args.GetSymbol(),
	}

	response, err := c.api.RequestWithdrawal(&request)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %w", err)
	}

	status := client.OperationStatusPending
	if response.Status == "confirmed" {
		status = client.OperationStatusSuccess
	}

	return &client.WithdrawalResponse{
		ID:     fmt.Sprintf("%d", response.ID),
		Status: status,
	}, nil
}

func (c *Client) GetDepositAddress(args client.GetDepositAddressArgs) (oc.Address, error) {
	request := api.DepositAddressRequest{
		Blockchain: args.GetNetwork(),
	}

	response, err := c.api.GetDepositAddress(&request)
	if err != nil {
		return "", fmt.Errorf("failed to get deposit address: %w", err)
	}

	return oc.Address(response.Address), nil
}

func (c *Client) ListWithdrawalHistory(args client.WithdrawalHistoryArgs) ([]*client.WithdrawalHistory, error) {
	// Prepare request parameters
	var request api.WithdrawalHistoryRequest

	// Apply filters from args if provided
	// if args.GetStartTime() != nil {
	// 	startTime := args.GetStartTime().UnixMilli()
	// 	request.From = &startTime
	// }

	// if args.GetEndTime() != nil {
	// 	endTime := args.GetEndTime().UnixMilli()
	// 	request.To = &endTime
	// }

	if args.GetLimit() > 0 {
		limit := uint64(args.GetLimit())
		request.Limit = &limit
	}

	// Get withdrawal history from API
	response, err := c.api.GetWithdrawals(&request)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}

	// Convert API response to client format
	history := make([]*client.WithdrawalHistory, 0, len(response))
	for _, withdrawal := range response {
		status := client.OperationStatusPending
		if withdrawal.Status == "confirmed" {
			status = client.OperationStatusSuccess
		}

		history = append(history, &client.WithdrawalHistory{
			ID:            fmt.Sprintf("%d", withdrawal.ID),
			Status:        status,
			Symbol:        withdrawal.Symbol,
			Network:       oc.NetworkId(withdrawal.Blockchain),
			Amount:        withdrawal.Quantity,
			Fee:           withdrawal.Fee,
			TransactionId: withdrawal.TransactionHash,
			Comment:       withdrawal.Status,
			Notes: map[string]string{
				"createdAt":  withdrawal.CreatedAt,
				"isInternal": strconv.FormatBool(withdrawal.IsInternal),
			},
		})
	}

	return history, nil
}
