package binanceus

import (
	"fmt"
	"time"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/binanceus/api"
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
	api := api.NewClient(apiKey, secretKey)
	return &Client{api: api}, nil
}

func (c *Client) ListAssets() ([]*oc.Asset, error) {
	response, err := c.api.GetAssetConfig(&api.GetAssetConfigRequest{
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	assets := []*oc.Asset{}
	for _, asset := range response {
		for _, network := range asset.NetworkList {
			assets = append(assets, &oc.Asset{
				SymbolId:  asset.Coin,
				NetworkId: network.Network,
				// unfortunately binanceus doesn't provide contract addresses
				ContractAddress: "",
			})
		}
	}
	return assets, nil
}

func (c *Client) ListBalances(args client.GetBalanceArgs) ([]*client.BalanceDetail, error) {
	response, err := c.api.GetAccount(&api.GetAccountRequest{
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	balances := []*client.BalanceDetail{}
	for _, asset := range response.Balances {
		if asset.Free.IsZero() && asset.Locked.IsZero() {
			continue
		}
		balances = append(balances, &client.BalanceDetail{
			SymbolId:    asset.Asset,
			Available:   asset.Free,
			Unavailable: asset.Locked,
		})
	}
	return balances, nil
}

func (c *Client) CreateAccountTransfer(args client.AccountTransferArgs) (*client.TransferStatus, error) {
	// response, err := c.api.SubAccountTransfer(&api.SubAccountTransferRequest{
	// 	FromEmail:       args.GetFrom().Id(),
	// 	ToEmail:         args.GetTo().Id(),
	// 	Asset:           args.GetSymbol(),
	// 	Amount:          args.GetAmount(),
	// 	TimestampMillis: time.Now().UnixMilli(),
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create account transfer: %w", err)
	// }
	// return &client.TransferStatus{
	// 	ID:     response.TxnId,
	// 	Status: client.OperationStatusSuccess,
	// }, nil
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) CreateWithdrawal(args client.WithdrawalArgs) (*client.WithdrawalResponse, error) {
	response, err := c.api.WithdrawCrypto(&api.WithdrawRequest{
		Coin:            args.GetSymbol(),
		Network:         args.GetNetwork(),
		Address:         args.GetAddress(),
		Amount:          args.GetAmount(),
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %w", err)
	}
	return &client.WithdrawalResponse{
		ID:     response.Id,
		Status: client.OperationStatusPending,
	}, nil
}

func (c *Client) GetDepositAddress(args client.GetDepositAddressArgs) (oc.Address, error) {
	response, err := c.api.GetDepositAddress(&api.GetDepositAddressRequest{
		Coin:            args.GetSymbol(),
		Network:         args.GetNetwork(),
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get deposit address: %w", err)
	}
	return response.Address, nil
}

func (c *Client) ListWithdrawalHistory(args client.WithdrawalHistoryArgs) ([]*client.WithdrawalHistory, error) {
	response, err := c.api.GetCryptoWithdrawalHistory(&api.GetCryptoWithdrawalHistoryRequest{
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}
	history := []*client.WithdrawalHistory{}
	for _, record := range response {
		status := client.OperationStatusPending
		switch record.Status {
		case api.WithdrawalStatusCompleted:
			status = client.OperationStatusSuccess
		case api.WithdrawalStatusCanceled, api.WithdrawalStatusRejected, api.WithdrawalStatusFailure:
			status = client.OperationStatusFailed
		}
		history = append(history, &client.WithdrawalHistory{
			ID:            record.Id,
			Status:        status,
			Amount:        record.Amount,
			Symbol:        record.Coin,
			Network:       record.Network,
			Fee:           record.TransactionFee,
			TransactionId: record.TxId,
			Comment:       record.Status.String(),
			Notes:         map[string]string{},
		})
	}
	return history, nil
}
