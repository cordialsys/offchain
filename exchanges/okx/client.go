package okx

import (
	"fmt"

	rebal "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/okx/api"
	"github.com/cordialsys/offchain/pkg/debug"
)

type Client struct {
	api *api.Client
}

var _ rebal.Client = &Client{}

func NewClient(config *rebal.ExchangeConfig) (*Client, error) {
	api, err := api.NewClient(config.ApiKey, config.SecretKey, config.Passphrase)
	if err != nil {
		return nil, err
	}
	return &Client{
		api: api,
	}, nil
}
func (c *Client) Api() *api.Client {
	return c.api
}

func (c *Client) GetAssets() ([]*rebal.Asset, error) {
	response, err := c.api.GetCurrencies()
	if err != nil {
		return nil, err
	}
	assets := make([]*rebal.Asset, len(response.Data))
	for i, currency := range response.Data {
		assets[i] = rebal.NewAsset(
			rebal.SymbolId(currency.Currency),
			rebal.NetworkId(currency.Chain),
			currency.ContractAddress,
		)
	}
	return assets, nil
}

func (c *Client) GetBalances(args rebal.GetBalanceArgs) ([]*rebal.BalanceDetail, error) {
	response, err := c.api.GetBalance(args)
	if err != nil {
		return nil, err
	}
	balances := []*rebal.BalanceDetail{}
	for _, balance := range response.Data {
		balances = append(balances, &rebal.BalanceDetail{
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

func (c *Client) AccountTransfer(args *rebal.AccountTransferArgs) error {

	from := ""
	to := ""
	switch account := args.GetFrom(); account {
	case rebal.CoreFunding:
		from = FundingAcount
	case rebal.CoreTrading:
		from = TradingAcount
	default:
		return fmt.Errorf("unsupported from account type: %s", account)
	}
	switch account := args.GetTo(); account {
	case rebal.CoreFunding:
		to = FundingAcount
	case rebal.CoreTrading:
		to = TradingAcount
	default:
		return fmt.Errorf("unsupported to account type: %s", account)
	}

	response, err := c.api.FundsTransfer(&api.AccountTransferRequest{
		From:     from,
		To:       to,
		Currency: args.GetSymbol(),
		Amount:   args.GetAmount(),
	})
	if err != nil {
		return err
	}
	debug.PrintJson(response)
	return nil
}
