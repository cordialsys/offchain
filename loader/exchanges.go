package loader

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/exchanges/backpack"
	"github.com/cordialsys/offchain/exchanges/binance"
	"github.com/cordialsys/offchain/exchanges/binanceus"
	"github.com/cordialsys/offchain/exchanges/bybit"
	"github.com/cordialsys/offchain/exchanges/okx"
)

type Client interface {
	client.Client
	// These methods basically just read the config
	ListAccountTypes() ([]*oc.AccountTypeConfig, error)
	ListSubaccounts() ([]*oc.SubAccountHeader, error)
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

func (c *ClientExtra) ListSubaccounts() ([]*oc.SubAccountHeader, error) {
	subaccounts := []*oc.SubAccountHeader{}
	for _, subaccount := range c.cfg.SubAccounts {
		subaccounts = append(subaccounts, &subaccount.SubAccountHeader)
	}
	return subaccounts, nil
}

func NewClient(config *oc.ExchangeConfig, account *oc.Account) (Client, error) {
	var cli client.Client
	var err error
	switch config.ExchangeId {
	case oc.Okx:
		cli, err = okx.NewClient(&config.ExchangeClientConfig, account)
	case oc.Bybit:
		cli, err = bybit.NewClient(&config.ExchangeClientConfig, account)
	case oc.Binance:
		cli, err = binance.NewClient(&config.ExchangeClientConfig, account)
	case oc.BinanceUS:
		cli, err = binanceus.NewClient(&config.ExchangeClientConfig, account)
	case oc.Backpack:
		cli, err = backpack.NewClient(&config.ExchangeClientConfig, account)
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", config.ExchangeId)
	}
	if err != nil {
		return nil, err
	}
	return newClient(config, cli), nil
}

func LoadValidatedConfig(path string) (*oc.Config, error) {
	cfg, err := oc.LoadUnvalidatedConfig(path)
	if err != nil {
		return nil, err
	}
	// general validation
	for _, exchange := range cfg.Exchanges {
		subaccountAliases := map[string]bool{}
		for _, subaccount := range exchange.SubAccounts {
			if subaccount.Alias == "" {
				continue
			}
			if _, ok := subaccountAliases[subaccount.Alias]; ok {
				return nil, fmt.Errorf("subaccount alias %s for exchange %s is duplicated", subaccount.Alias, exchange.ExchangeId)
			}
			subaccountAliases[subaccount.Alias] = true
		}
	}
	// per-exchange validation
	for _, exchange := range cfg.Exchanges {
		switch exchange.ExchangeId {
		case oc.Bybit:
			if err := bybit.Validate(exchange); err != nil {
				return nil, err
			}
		case oc.Binance:
			if err := binance.Validate(exchange); err != nil {
				return nil, err
			}
		case oc.BinanceUS:
			if err := binanceus.Validate(exchange); err != nil {
				return nil, err
			}
		case oc.Okx:
			if err := okx.Validate(exchange); err != nil {
				return nil, err
			}
		case oc.Backpack:
			if err := backpack.Validate(exchange); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported exchange in config: %s", exchange.ExchangeId)
		}
	}
	return cfg, nil
}
