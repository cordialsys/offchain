package offchain

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Exchanges map[ExchangeId]*ExchangeConfig `yaml:"exchanges"`
}

func (cfg *Config) Init() error {
	for key, exchange := range cfg.Exchanges {
		if !slices.Contains(ValidExchangeIds, key) {
			return fmt.Errorf("invalid exchange id: %s", key)
		}
		if exchange == nil {
			return fmt.Errorf("empty config for exchange: %s", key)
		}
		exchange.ExchangeId = key
	}
	return nil
}

func (c *Config) GetExchange(id ExchangeId) (*ExchangeConfig, bool) {
	exchange, ok := c.Exchanges[id]
	if !ok {
		return nil, false
	}
	return exchange, true
}

const ENV_OFFCHAIN_CONFIG = "OFFCHAIN_CONFIG"

func LoadUnvalidatedConfig(configPathMaybe string) (*Config, error) {
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

	err = section.Offchain.Init()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	for key, exchange := range section.Offchain.Exchanges {

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

	for _, exchange := range section.Offchain.Exchanges {
		// apply defaults to fields not overridden by user
		ApplyDefaults(exchange)
	}

	return &section.Offchain, nil
}
