package loader

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/binance"
	"github.com/cordialsys/offchain/exchanges/binanceus"
	"github.com/cordialsys/offchain/exchanges/bybit"
	"github.com/cordialsys/offchain/exchanges/okx"
	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func NewClient(config *oc.ExchangeConfig) (oc.Client, error) {
	switch config.ExchangeId {
	case oc.Okx:
		return okx.NewClient(config)
	case oc.Bybit:
		return bybit.NewClient(config)
	case oc.Binance:
		return binance.NewClient(config)
	case oc.BinanceUS:
		return binanceus.NewClient(config)
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", config.ExchangeId)
	}
}

type Config struct {
	Exchanges map[oc.ExchangeId]*oc.ExchangeConfig `yaml:"exchanges"`
}

func (c *Config) GetExchange(id oc.ExchangeId) (*oc.ExchangeConfig, bool) {
	exchange, ok := c.Exchanges[id]
	if !ok {
		return nil, false
	}
	return exchange, true
}

const ENV_OFFCHAIN_CONFIG = "OFFCHAIN_CONFIG"

func LoadConfig(configPathMaybe string) (*Config, error) {
	type configsection struct {
		Offchain Config `yaml:"offchain"`
	}
	section := configsection{}
	if configPathMaybe == "" {
		if v := os.Getenv(ENV_OFFCHAIN_CONFIG); v != "" {
			configPathMaybe = v
		}
	}
	if configPathMaybe == "" {
		return nil, fmt.Errorf("config path is required (maybe set %s)", ENV_OFFCHAIN_CONFIG)
	}

	logrus.WithField("config", configPathMaybe).Info("loading configuration")
	err := cleanenv.ReadConfig(configPathMaybe, &section)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration: %v", err)
	}

	for key, exchange := range section.Offchain.Exchanges {
		if !slices.Contains(oc.ValidExchangeIds, key) {
			return nil, fmt.Errorf("invalid exchange id: %s", key)
		}
		exchange.ExchangeId = key
		if exchange.ApiKeyRef.IsType(secret.Raw) {
			slog.Warn("raw api-key in config file is unsafe, should use a secret manager", "exchange", key)
		}
		if exchange.SecretKeyRef.IsType(secret.Raw) {
			slog.Warn("raw secret-key in config file is unsafe, should use a secret manager", "exchange", key)
		}
		if exchange.PassphraseRef.IsType(secret.Raw) {
			slog.Warn("raw passphrase in config file is unsafe, should use a secret manager", "exchange", key)
		}
	}

	return &section.Offchain, nil
}
