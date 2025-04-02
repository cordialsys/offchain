package binance

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/binance/api"
	"github.com/cordialsys/offchain/pkg/debug"
)

type Client struct {
	api *api.Client
}

var _ client.Client = &Client{}

func NewClient(config *oc.ExchangeConfig) (*Client, error) {
	apiKey, err := config.ApiKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load api key: %w", err)
	}
	secretKey, err := config.SecretKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load secret key: %w", err)
	}
	api, err := api.NewClient(apiKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create api client: %w", err)
	}
	return &Client{api: api}, nil
}

func (c *Client) ListAssets() ([]*oc.Asset, error) {
	response, err := c.api.GetAllCoinsInformation()
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	debug.PrintJson(response)
	assets := []*oc.Asset{}
	for _, asset := range response {
		for _, network := range asset.NetworkList {
			assets = append(assets, &oc.Asset{
				SymbolId:        asset.Coin,
				NetworkId:       network.Network,
				ContractAddress: network.ContractAddress,
			})
		}
	}
	return assets, nil
}

func (c *Client) ListBalances(args client.GetBalanceArgs) ([]*client.BalanceDetail, error) {
	response, err := c.api.GetUserAsset(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	balances := make([]*client.BalanceDetail, len(response))
	for i, asset := range response {
		// not sure what the difference between frozen and locked is for binance
		frozen := asset.Freeze.Decimal()
		locked := asset.Locked.Decimal()
		unavailable := frozen.Add(locked)

		balances[i] = &client.BalanceDetail{
			SymbolId:    asset.Asset,
			Available:   asset.Free,
			Unavailable: oc.Amount(unavailable),
		}
	}
	return balances, nil
}

// probably should use this: https://developers.binance.com/docs/sub_account/asset-management/Universal-Transfer
// TODO
func (c *Client) CreateAccountTransfer(args client.AccountTransferArgs) (*client.TransferStatus, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) CreateWithdrawal(args client.WithdrawalArgs) (*client.WithdrawalResponse, error) {
	req := api.WithdrawRequest{
		Coin:    args.GetSymbol(),
		Network: args.GetNetwork(),
		Address: args.GetAddress(),
		Amount:  args.GetAmount(),
	}

	response, err := c.api.Withdraw(&req)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %w", err)
	}
	return &client.WithdrawalResponse{
		ID:     response.Id,
		Status: client.OperationStatusPending,
	}, nil
}

func (c *Client) GetDepositAddress(args client.GetDepositAddressArgs) (oc.Address, error) {
	request := api.DepositAddressRequest{
		Coin:    args.GetSymbol(),
		Network: args.GetNetwork(),
		Amount:  nil,
	}

	response, err := c.api.GetDepositAddress(&request)
	if err != nil {
		return "", fmt.Errorf("failed to get deposit address: %w", err)
	}
	return response.Address, nil
}

func (c *Client) ListWithdrawalHistory(args client.WithdrawalHistoryArgs) ([]*client.WithdrawalHistory, error) {
	response, err := c.api.GetWithdrawalHistory(&api.WithdrawalHistoryRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}
	history := []*client.WithdrawalHistory{}
	for _, record := range response {
		status := client.OperationStatusPending
		switch record.Status {
		case api.WithdrawalStatusCompleted:
			status = client.OperationStatusSuccess
		case api.WithdrawalStatusProcessing, api.WithdrawalStatusAwaitingApproval, api.WithdrawalStatusEmailSent:
			status = client.OperationStatusPending
		case api.WithdrawalStatusRejected:
			status = client.OperationStatusFailed
		}

		history = append(history, &client.WithdrawalHistory{
			ID:            record.Id,
			Status:        status,
			Symbol:        record.Coin,
			Network:       record.Network,
			Amount:        record.Amount,
			Fee:           record.TransactionFee,
			TransactionId: record.TxId,
			Comment:       record.Info,
			Notes:         map[string]string{},
		})
	}
	return history, nil
}
