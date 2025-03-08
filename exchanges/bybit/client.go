package bybit

import (
	"fmt"
	"log/slog"
	"time"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/bybit/api"
)

type Client struct {
	api *api.Client
}

var _ oc.Client = &Client{}

func NewClient(config *oc.ExchangeConfig) (*Client, error) {
	apiKey, err := config.ApiKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load api key: %v", err)
	}
	secretKey, err := config.SecretKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load secret key: %v", err)
	}

	api, err := api.NewClient(apiKey, secretKey)
	if err != nil {
		return nil, err
	}
	if config.ApiUrl != "" {
		api.SetBaseURL(config.ApiUrl)
	}
	return &Client{
		api: api,
	}, nil
}

func (c *Client) ListAssets() ([]*oc.Asset, error) {
	response, err := c.api.GetCoinInfo("")
	if err != nil {
		return nil, err
	}
	assets := []*oc.Asset{}
	for _, coin := range response.Result.Rows {
		for _, chain := range coin.Chains {
			assets = append(assets, &oc.Asset{
				SymbolId:        oc.SymbolId(coin.Coin),
				NetworkId:       oc.NetworkId(chain.Chain),
				ContractAddress: oc.ContractAddress(chain.ContractAddress),
			})
		}
	}
	return assets, nil
}

func mapAccountName(name oc.AccountName) string {
	switch name {
	case "":
		// default to funding
		return api.AccountTypeFund
	case oc.CoreFunding:
		return api.AccountTypeFund
	case oc.CoreTrading:
		return api.AccountTypeUnified
	default:
		return name.Id()
	}
}

func (c *Client) ListBalances(args oc.GetBalanceArgs) ([]*oc.BalanceDetail, error) {

	accountType := mapAccountName(args.GetAccount())

	response, err := c.api.GetAllCoinsBalance(accountType, "")
	if err != nil {
		return nil, err
	}
	balances := []*oc.BalanceDetail{}
	for _, balance := range response.Result.Balance {

		walletBalance := balance.WalletBalance.Decimal()
		transferBalance := balance.TransferBalance.Decimal()
		frozen := walletBalance.Sub(transferBalance)

		balances = append(balances, &oc.BalanceDetail{
			Available:   oc.Amount(transferBalance),
			Unavailable: oc.Amount(frozen),
			SymbolId:    balance.Coin,
			NetworkId:   "",
		})
	}

	return balances, nil
}

func (c *Client) CreateAccountTransfer(args oc.AccountTransferArgs) (*oc.TransferStatus, error) {
	from := mapAccountName(args.GetFrom())
	to := mapAccountName(args.GetTo())
	coin := args.GetSymbol()
	amount := args.GetAmount()

	if amount.IsZero() {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	response, err := c.api.CreateInternalTransfer(coin, amount, from, to)
	if err != nil {
		return nil, err
	}
	state := oc.OperationStatusSuccess
	if response.Result.Status != api.TransferStateSuccess {
		slog.Error("transfer status unexpected", "status", response.Result.Status)
		// assume pending as it could still occur
		state = oc.OperationStatusPending
	}
	return &oc.TransferStatus{
		ID:     response.Result.TransferID,
		Status: state,
	}, nil
}

func (c *Client) CreateWithdrawal(args oc.WithdrawalArgs) (*oc.WithdrawalResponse, error) {

	request := api.WithdrawRequest{
		Coin:        args.GetSymbol(),
		Address:     args.GetAddress(),
		Chain:       args.GetNetwork(),
		Amount:      args.GetAmount(),
		Timestamp:   time.Now().UnixMilli(),
		AccountType: api.AccountTypeFund,
	}

	response, err := c.api.Withdraw(&request)
	if err != nil {
		return nil, err
	}
	return &oc.WithdrawalResponse{
		ID:     response.Result.ID,
		Status: oc.OperationStatusPending,
	}, nil
}

func (c *Client) GetDepositAddress(args oc.GetDepositAddressArgs) (oc.Address, error) {
	var response *api.GetDepositAddressResponse
	var err error
	if accountId, ok := args.GetAccountId(); ok {
		response, err = c.api.GetSubDepositAddress(accountId, args.GetSymbol(), args.GetNetwork())
	} else {
		response, err = c.api.GetMasterDepositAddress(args.GetSymbol(), args.GetNetwork())
	}

	if err != nil {
		return "", err
	}
	for _, chain := range response.Result.Chains {
		if chain.Chain == args.GetNetwork() {
			return oc.Address(chain.AddressDeposit), nil
		}
	}
	return "", fmt.Errorf("no deposit address found for network %s", args.GetNetwork())
}

func (c *Client) ListWithdrawalHistory(args oc.WithdrawalHistoryArgs) ([]*oc.WithdrawalHistory, error) {
	response, err := c.api.GetWithdrawalRecords(&api.WithdrawalRecordsRequest{})
	if err != nil {
		return nil, err
	}
	history := []*oc.WithdrawalHistory{}
	for _, record := range response.Result.Rows {

		status := oc.OperationStatusPending
		switch record.Status {
		case api.WithdrawalStatusSuccess:
			status = oc.OperationStatusSuccess
		case api.WithdrawalStatusFail,
			api.WithdrawalStatusCancelByUser,
			api.WithdrawalStatusMoreInformationRequired,
			api.WithdrawalStatusReject:
			status = oc.OperationStatusFailed
		}

		history = append(history, &oc.WithdrawalHistory{
			ID:            record.WithdrawId,
			Status:        status,
			Symbol:        record.Coin,
			Network:       record.Chain,
			Amount:        record.Amount,
			Fee:           record.WithdrawFee,
			TransactionId: record.TxID,
			Comment:       string(record.Status),
			Notes:         map[string]string{},
		})
	}
	return history, nil
}
