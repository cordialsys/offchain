package okx

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/okx/api"
)

type Client struct {
	api *api.Client
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
	passphrase, err := secrets.PassphraseRef.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load passphrase: %v", err)
	}
	api, err := api.NewClient(apiKey, secretKey, passphrase)
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
func (c *Client) Api() *api.Client {
	return c.api
}

func (c *Client) ListAssets() ([]*oc.Asset, error) {
	response, err := c.api.GetCurrencies()
	if err != nil {
		return nil, err
	}
	assets := make([]*oc.Asset, len(response.Data))
	for i, currency := range response.Data {
		assets[i] = oc.NewAsset(
			oc.SymbolId(currency.Currency),
			oc.NetworkId(currency.Chain),
			currency.ContractAddress,
		)
	}
	return assets, nil
}

func (c *Client) ListBalances(args client.GetBalanceArgs) ([]*client.BalanceDetail, error) {
	response, err := c.api.GetBalance(args)
	if err != nil {
		return nil, err
	}
	balances := []*client.BalanceDetail{}
	for _, balance := range response.Data {
		balances = append(balances, &client.BalanceDetail{
			SymbolId: balance.Currency,
			// NA
			NetworkId:   "",
			Available:   balance.AvailableBalance,
			Unavailable: balance.FrozenBalance,
		})
	}
	return balances, nil
}

// Not sure why but OKX use these static numbers to refer to the core accounts
const FundingAcount = "6"
const TradingAcount = "18"

func (c *Client) CreateAccountTransfer(args client.AccountTransferArgs) (*client.TransferStatus, error) {

	// from := ""
	// to := ""
	// switch account := args.GetFrom(); account {
	// case client.CoreFunding:
	// 	from = FundingAcount
	// case client.CoreTrading:
	// 	from = TradingAcount
	// default:
	// 	return nil, fmt.Errorf("unsupported from account type: %s", account)
	// }
	// switch account := args.GetTo(); account {
	// case client.CoreFunding:
	// 	to = FundingAcount
	// case client.CoreTrading:
	// 	to = TradingAcount
	// default:
	// 	return nil, fmt.Errorf("unsupported to account type: %s", account)
	// }

	// response, err := c.api.FundsTransfer(&api.AccountTransferRequest{
	// 	From:     from,
	// 	To:       to,
	// 	Currency: args.GetSymbol(),
	// 	Amount:   args.GetAmount(),
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return &client.TransferStatus{
	// 	ID:     response.Data[0].TransferID,
	// 	Status: client.OperationStatusSuccess,
	// }, nil
	return nil, fmt.Errorf("not implemented")
}

// more special static numbers okx used
const InternalTransfer = "3"
const OnChainTransfer = "4"

func (c *Client) CreateWithdrawal(args client.WithdrawalArgs) (*client.WithdrawalResponse, error) {
	response, err := c.api.Withdrawal(&api.WithdrawalRequest{
		Amount:         args.GetAmount().String(),
		Destination:    OnChainTransfer,
		Currency:       args.GetSymbol(),
		SymbolAndChain: fmt.Sprintf("%s-%s", args.GetSymbol(), args.GetNetwork()),
		ToAddress:      args.GetAddress(),
	})
	if err != nil {
		return nil, err
	}
	return &client.WithdrawalResponse{
		ID:     response.Data[0].WithdrawalID,
		Status: client.OperationStatusPending,
	}, nil
}

func (c *Client) GetDepositAddress(args client.GetDepositAddressArgs) (oc.Address, error) {
	response, err := c.api.GetDepositAddress(args.GetSymbol())
	if err != nil {
		return "", err
	}
	expectedSymbolAndChain := fmt.Sprintf("%s-%s", args.GetSymbol(), args.GetNetwork())
	for _, chain := range response.Data {
		if chain.SymbolAndChain == expectedSymbolAndChain {
			return oc.Address(chain.Address), nil
		}
	}
	return "", fmt.Errorf("no deposit address found for network %s", args.GetNetwork())
}

func (c *Client) ListWithdrawalHistory(args client.WithdrawalHistoryArgs) ([]*client.WithdrawalHistory, error) {
	response, err := c.api.GetWithdrawalHistory(&api.WithdrawalHistoryRequest{})
	if err != nil {
		return nil, err
	}
	history := []*client.WithdrawalHistory{}
	for _, record := range response.Data {
		status := client.OperationStatusPending
		if record.State == api.WithdrawalStateSuccess {
			status = client.OperationStatusSuccess
		} else if record.State == api.WithdrawalStateFailed || record.State == api.WithdrawalStateCanceled {
			status = client.OperationStatusFailed
		}
		withdrawal := &client.WithdrawalHistory{
			ID:            record.WithdrawalId,
			Status:        status,
			Symbol:        record.Currency,
			Network:       record.Chain.NetworkId(),
			Amount:        record.Amount,
			Fee:           record.Fee,
			TransactionId: record.TxId,
			Comment:       record.Note,
			Notes:         map[string]string{},
		}

		history = append(history, withdrawal)
	}
	return history, nil
}
