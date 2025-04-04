package bybit

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/bybit/api"
)

type Client struct {
	api *api.Client
	cfg *oc.ExchangeClientConfig
}

var _ client.Client = &Client{}

func NewClient(config *oc.ExchangeClientConfig, secrets *oc.MultiSecret) (*Client, error) {
	apiKey, err := secrets.ApiKeyRef.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load api key: %v", err)
	}
	secretKey, err := secrets.SecretKeyRef.Load()
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
		cfg: config,
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

func (c *Client) ListBalances(args client.GetBalanceArgs) ([]*client.BalanceDetail, error) {

	accountType := args.GetAccountType()

	response, err := c.api.GetAllCoinsBalance(accountType, "")
	if err != nil {
		return nil, err
	}
	balances := []*client.BalanceDetail{}
	for _, balance := range response.Result.Balance {

		walletBalance := balance.WalletBalance.Decimal()
		transferBalance := balance.TransferBalance.Decimal()
		frozen := walletBalance.Sub(transferBalance)

		balances = append(balances, &client.BalanceDetail{
			Available:   oc.Amount(transferBalance),
			Unavailable: oc.Amount(frozen),
			SymbolId:    balance.Coin,
			NetworkId:   "",
		})
	}

	return balances, nil
}

func (c *Client) CreateAccountTransfer(args client.AccountTransferArgs) (*client.TransferStatus, error) {
	var err error
	var response *api.TransferResponse

	if args.IsSameAccount() {
		// apparently the 'universal' transfer doesn't work for transfers in the same account
		_, fromType := args.GetFrom()
		_, toType := args.GetTo()
		response, err = c.api.CreateInternalTransfer(args.GetSymbol(), args.GetAmount(), fromType, toType)
		if err != nil {
			return nil, err
		}
	} else {
		from, fromType := args.GetFrom()
		to, toType := args.GetTo()

		if from == "" {
			// use the main account
			from = c.cfg.Id
		}
		if to == "" {
			// use the main account
			to = c.cfg.Id
		}

		fromId, err := strconv.Atoi(string(from))
		if err != nil {
			return nil, fmt.Errorf("expected numeric UID for from-account, not %s: %v", from, err)
		}
		toId, err := strconv.Atoi(string(to))
		if err != nil {
			return nil, fmt.Errorf("expected numeric UID for to-account, not %s: %v", to, err)
		}

		response, err = c.api.CreateUniversalTransfer(
			args.GetSymbol(),
			args.GetAmount(),
			fromId,
			toId,
			fromType,
			toType,
		)
		if err != nil {
			return nil, err
		}
	}

	state := client.OperationStatusSuccess
	if response.Result.Status != api.TransferStateSuccess {
		slog.Error("transfer status unexpected", "status", response.Result.Status)
		// assume pending as it could still occur
		state = client.OperationStatusPending
	}
	return &client.TransferStatus{
		ID:     response.Result.TransferID,
		Status: state,
	}, nil
}

func (c *Client) CreateWithdrawal(args client.WithdrawalArgs) (*client.WithdrawalResponse, error) {

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
	return &client.WithdrawalResponse{
		ID:     response.Result.ID,
		Status: client.OperationStatusPending,
	}, nil
}

func (c *Client) GetDepositAddress(args client.GetDepositAddressArgs) (oc.Address, error) {
	var response *api.GetDepositAddressResponse
	var err error
	if accountType, ok := args.GetSubaccount(); ok {
		response, err = c.api.GetSubDepositAddress(accountType, args.GetSymbol(), args.GetNetwork())
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

func (c *Client) ListWithdrawalHistory(args client.WithdrawalHistoryArgs) ([]*client.WithdrawalHistory, error) {
	response, err := c.api.GetWithdrawalRecords(&api.WithdrawalRecordsRequest{})
	if err != nil {
		return nil, err
	}
	history := []*client.WithdrawalHistory{}
	for _, record := range response.Result.Rows {

		status := client.OperationStatusPending
		switch record.Status {
		case api.WithdrawalStatusSuccess:
			status = client.OperationStatusSuccess
		case api.WithdrawalStatusFail,
			api.WithdrawalStatusCancelByUser,
			api.WithdrawalStatusMoreInformationRequired,
			api.WithdrawalStatusReject:
			status = client.OperationStatusFailed
		}

		history = append(history, &client.WithdrawalHistory{
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
