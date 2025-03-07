package okx

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/okx/api"
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
	passphrase, err := config.PassphraseRef.Load()
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

func (c *Client) ListBalances(args oc.GetBalanceArgs) ([]*oc.BalanceDetail, error) {
	response, err := c.api.GetBalance(args)
	if err != nil {
		return nil, err
	}
	balances := []*oc.BalanceDetail{}
	for _, balance := range response.Data {
		balances = append(balances, &oc.BalanceDetail{
			SymbolId: balance.Currency,
			// NA
			NetworkId: "",
			Available: balance.AvailableBalance,
			Frozen:    balance.FrozenBalance,
		})
	}
	return balances, nil
}

// Not sure why but OKX use these static numbers to refer to the core accounts
const FundingAcount = "6"
const TradingAcount = "18"

func (c *Client) CreateAccountTransfer(args *oc.AccountTransferArgs) (*oc.TransferStatus, error) {

	from := ""
	to := ""
	switch account := args.GetFrom(); account {
	case oc.CoreFunding:
		from = FundingAcount
	case oc.CoreTrading:
		from = TradingAcount
	default:
		return nil, fmt.Errorf("unsupported from account type: %s", account)
	}
	switch account := args.GetTo(); account {
	case oc.CoreFunding:
		to = FundingAcount
	case oc.CoreTrading:
		to = TradingAcount
	default:
		return nil, fmt.Errorf("unsupported to account type: %s", account)
	}

	response, err := c.api.FundsTransfer(&api.AccountTransferRequest{
		From:     from,
		To:       to,
		Currency: args.GetSymbol(),
		Amount:   args.GetAmount(),
	})
	if err != nil {
		return nil, err
	}
	return &oc.TransferStatus{
		ID:     response.Data[0].TransferID,
		Status: oc.OperationStatusSuccess,
	}, nil
}

// more special static numbers okx used
const InternalTransfer = "3"
const OnChainTransfer = "4"

func (c *Client) CreateWithdrawal(args *oc.WithdrawalArgs) (*oc.WithdrawalResponse, error) {
	response, err := c.api.Withdrawal(&api.WithdrawalRequest{
		Amount:      args.GetAmount().String(),
		Destination: OnChainTransfer,
		Currency:    args.GetSymbol(),
		Chain:       args.GetNetwork(),
		ToAddress:   args.GetAddress(),
	})
	if err != nil {
		return nil, err
	}
	return &oc.WithdrawalResponse{
		ID:     response.Data[0].WithdrawalID,
		Status: oc.OperationStatusPending,
	}, nil
}
