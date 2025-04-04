package loader

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/binance"
	"github.com/cordialsys/offchain/exchanges/binanceus"
	"github.com/cordialsys/offchain/exchanges/bybit"
	"github.com/cordialsys/offchain/exchanges/okx"
)

type Client interface {
	client.Client
	// These methods basically just read the config
	ListAccountTypes() ([]*oc.AccountTypeConfig, error)
	ListSubaccounts() ([]string, error)
}

type ClientExtra struct {
	client.Client
	cfg *oc.ExchangeConfig
}

var _ Client = &ClientExtra{}

func newClient(cfg *oc.ExchangeConfig, client client.Client) Client {
	return &ClientExtra{
		Client: client,
		cfg:    cfg,
	}
}

func (c *ClientExtra) ListAccountTypes() ([]*oc.AccountTypeConfig, error) {
	return c.cfg.AccountTypes, nil
}

func (c *ClientExtra) ListSubaccounts() ([]string, error) {
	subaccounts := []string{}
	for _, subaccount := range c.cfg.SubAccounts {
		subaccounts = append(subaccounts, subaccount.Id)
	}
	return subaccounts, nil
}

func NewClient(config *oc.ExchangeConfig, secrets *oc.MultiSecret) (Client, error) {
	var cli client.Client
	var err error
	switch config.ExchangeId {
	case oc.Okx:
		cli, err = okx.NewClient(&config.ExchangeClientConfig, secrets)
	case oc.Bybit:
		cli, err = bybit.NewClient(&config.ExchangeClientConfig, secrets)
	case oc.Binance:
		cli, err = binance.NewClient(&config.ExchangeClientConfig, secrets)
	case oc.BinanceUS:
		cli, err = binanceus.NewClient(&config.ExchangeClientConfig, secrets)
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", config.ExchangeId)
	}
	if err != nil {
		return nil, err
	}
	return newClient(config, cli), nil
}
