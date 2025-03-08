package binanceus

import (
	"fmt"
	"time"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/binanceus/api"
)

type Client struct {
	api *api.Client
}

var _ oc.Client = &Client{}

func NewClient(config *oc.ExchangeConfig) (*Client, error) {
	apiKey, err := config.ApiKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load api key: %w", err)
	}
	secretKey, err := config.SecretKeyRef.Load()
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

func (c *Client) ListBalances(args oc.GetBalanceArgs) ([]*oc.BalanceDetail, error) {
	response, err := c.api.GetAccount(&api.GetAccountRequest{
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	balances := []*oc.BalanceDetail{}
	for _, asset := range response.Balances {
		if asset.Free.IsZero() && asset.Locked.IsZero() {
			continue
		}
		balances = append(balances, &oc.BalanceDetail{
			SymbolId:    asset.Asset,
			Available:   asset.Free,
			Unavailable: asset.Locked,
		})
	}
	return balances, nil
}

func (c *Client) CreateAccountTransfer(args oc.AccountTransferArgs) (*oc.TransferStatus, error) {
	response, err := c.api.SubAccountTransfer(&api.SubAccountTransferRequest{
		FromEmail:       args.GetFrom().Id(),
		ToEmail:         args.GetTo().Id(),
		Asset:           args.GetSymbol(),
		Amount:          args.GetAmount(),
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create account transfer: %w", err)
	}
	return &oc.TransferStatus{
		ID:     response.TxnId,
		Status: oc.OperationStatusSuccess,
	}, nil
}

func (c *Client) CreateWithdrawal(args oc.WithdrawalArgs) (*oc.WithdrawalResponse, error) {
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
	return &oc.WithdrawalResponse{
		ID:     response.Id,
		Status: oc.OperationStatusPending,
	}, nil
}

func (c *Client) GetDepositAddress(args oc.GetDepositAddressArgs) (oc.Address, error) {
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

func (c *Client) ListWithdrawalHistory(args oc.WithdrawalHistoryArgs) ([]*oc.WithdrawalHistory, error) {
	response, err := c.api.GetCryptoWithdrawalHistory(&api.GetCryptoWithdrawalHistoryRequest{
		TimestampMillis: time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}
	history := []*oc.WithdrawalHistory{}
	for _, record := range response {
		status := oc.OperationStatusPending
		switch record.Status {
		case api.WithdrawalStatusCompleted:
			status = oc.OperationStatusSuccess
		case api.WithdrawalStatusCanceled, api.WithdrawalStatusRejected, api.WithdrawalStatusFailure:
			status = oc.OperationStatusFailed
		}
		history = append(history, &oc.WithdrawalHistory{
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
