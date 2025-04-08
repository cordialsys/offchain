package offchain

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	serverclient "github.com/cordialsys/offchain/server/client"
	"github.com/cordialsys/offchain/server/client/api"
)

// This is a "simulated" exchange that just proxies requests to an offchain server
type Client struct {
	cli      *serverclient.Client
	exchange oc.ExchangeId
}

var _ client.Client = &Client{}

func NewClient(serverClient *serverclient.Client, exchange oc.ExchangeId) (*Client, error) {
	return &Client{
		cli:      serverClient,
		exchange: exchange,
	}, nil
}

func (c *Client) ListAssets() ([]*oc.Asset, error) {
	assets, err := c.cli.GetAssets(c.exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}

	// Convert from API assets to domain assets
	result := make([]*oc.Asset, len(assets))
	for i, asset := range assets {
		result[i] = &oc.Asset{
			SymbolId:        oc.SymbolId(asset.Symbol),
			NetworkId:       oc.NetworkId(asset.Network),
			ContractAddress: oc.ContractAddress(api.DerefOrZero(asset.Contract)),
		}
		if asset.Name != nil {
			// Handle any additional asset properties if needed
		}
	}

	return result, nil
}

func (c *Client) ListBalances(args client.GetBalanceArgs) ([]*client.BalanceDetail, error) {
	// Get account type from args if specified
	accountType := ""
	if args.GetAccountType() != "" {
		accountType = string(args.GetAccountType())
	}

	balances, err := c.cli.GetBalances(c.exchange, accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to get balances: %w", err)
	}

	// Convert from API balances to domain balances
	result := make([]*client.BalanceDetail, len(balances))
	for i, balance := range balances {
		available, err := oc.NewAmountFromString(balance.Available)
		if err != nil {
			return nil, fmt.Errorf("invalid available amount: %w", err)
		}

		unavailable, err := oc.NewAmountFromString(balance.Unavailable)
		if err != nil {
			return nil, fmt.Errorf("invalid unavailable amount: %w", err)
		}

		result[i] = &client.BalanceDetail{
			SymbolId:    oc.SymbolId(balance.Symbol),
			NetworkId:   oc.NetworkId(balance.Network),
			Available:   available,
			Unavailable: unavailable,
		}
	}

	return result, nil
}

func (c *Client) CreateAccountTransfer(args client.AccountTransferArgs) (*client.TransferStatus, error) {
	// Create transfer request
	transfer := api.NewTransfer(
		string(args.GetSymbol()),
		args.GetAmount(),
	)

	// Set from account if specified
	from, fromType := args.GetFrom()
	if from != "" {
		transfer = transfer.WithFrom(string(from))
	}
	if fromType != "" {
		transfer = transfer.WithFromType(string(fromType))
	}

	// Set to account if specified
	to, toType := args.GetTo()
	if to != "" {
		transfer = transfer.WithTo(string(to))
	}
	if toType != "" {
		transfer = transfer.WithToType(string(toType))
	}

	// Execute the transfer
	response, err := c.cli.CreateAccountTransfer(c.exchange, transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to create account transfer: %w", err)
	}

	return &client.TransferStatus{
		ID:     response.Id,
		Status: client.OperationStatus(response.Status),
	}, nil
}

func (c *Client) CreateWithdrawal(args client.WithdrawalArgs) (*client.WithdrawalResponse, error) {
	// Create withdrawal request
	withdrawal := api.NewWithdrawal(
		string(args.GetAddress()),
		string(args.GetSymbol()),
		string(args.GetNetwork()),
		args.GetAmount(),
	)

	// Execute the withdrawal
	response, err := c.cli.CreateWithdrawal(c.exchange, withdrawal)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %w", err)
	}

	return &client.WithdrawalResponse{
		ID:     response.Id,
		Status: client.OperationStatus(response.Status),
	}, nil
}

func (c *Client) GetDepositAddress(args client.GetDepositAddressArgs) (oc.Address, error) {

	forMaybe, _ := args.GetSubaccount()
	address, err := c.cli.GetDepositAddress(
		c.exchange,
		(args.GetSymbol()),
		(args.GetNetwork()),
		(forMaybe),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get deposit address: %w", err)
	}

	return oc.Address(address), nil
}

func (c *Client) ListWithdrawalHistory(args client.WithdrawalHistoryArgs) ([]*client.WithdrawalHistory, error) {
	limit := 0
	if args.GetLimit() > 0 {
		limit = args.GetLimit()
	}

	pageToken := ""
	if args.GetPageToken() != "" {
		pageToken = args.GetPageToken()
	}

	withdrawals, err := c.cli.ListWithdrawalHistory(c.exchange, limit, pageToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}

	// Convert from API withdrawals to domain withdrawals
	result := make([]*client.WithdrawalHistory, len(withdrawals))
	for i, withdrawal := range withdrawals {
		amount, err := oc.NewAmountFromString(withdrawal.Amount)
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}

		var fee oc.Amount
		if withdrawal.Fee != nil {
			feeAmount, err := oc.NewAmountFromString(*withdrawal.Fee)
			if err != nil {
				return nil, fmt.Errorf("invalid fee: %w", err)
			}
			fee = feeAmount
		}

		var txID client.TransactionId
		if withdrawal.TransactionId != nil {
			txID = client.TransactionId(*withdrawal.TransactionId)
		}

		notes := make(map[string]string)
		if withdrawal.Notes != nil {
			for k, v := range *withdrawal.Notes {
				notes[k] = v
			}
		}

		result[i] = &client.WithdrawalHistory{
			ID:            withdrawal.Id,
			Status:        client.OperationStatus(withdrawal.Status),
			Symbol:        oc.SymbolId(withdrawal.Symbol),
			Network:       oc.NetworkId(withdrawal.Network),
			Amount:        amount,
			Fee:           fee,
			TransactionId: txID,
			Comment:       api.DerefOrZero(withdrawal.Comment),
			Notes:         notes,
		}
	}

	return result, nil
}

func (c *Client) ListSubaccounts() ([]*oc.SubAccountHeader, error) {
	subaccounts, err := c.cli.ListSubaccounts(c.exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to list subaccounts: %w", err)
	}

	result := make([]*oc.SubAccountHeader, len(subaccounts))
	for i, subaccount := range subaccounts {
		result[i] = &oc.SubAccountHeader{
			Id: oc.AccountId(subaccount.Id),
		}
		if subaccount.Alias != nil {
			result[i].Alias = *subaccount.Alias
		}
	}

	return result, nil
}

func (c *Client) ListAccountTypes() ([]*oc.AccountTypeConfig, error) {
	accountTypes, err := c.cli.GetAccountTypes(c.exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to get account types: %w", err)
	}

	result := make([]*oc.AccountTypeConfig, len(accountTypes))
	for i, accountType := range accountTypes {
		result[i] = &oc.AccountTypeConfig{
			Type:    oc.AccountType(accountType.Type),
			Aliases: accountType.Aliases,
		}
		if accountType.Description != nil {
			result[i].Description = *accountType.Description
		}
	}

	return result, nil
}
