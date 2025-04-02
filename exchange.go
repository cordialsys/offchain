package offchain

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cordialsys/offchain/pkg/secret"
)

type ExchangeId string

var (
	Okx       ExchangeId = "okx"
	Binance   ExchangeId = "binance"
	BinanceUS ExchangeId = "binanceus"
	Bybit     ExchangeId = "bybit"
)

var ValidExchangeIds = []ExchangeId{Okx, Binance, BinanceUS, Bybit}

type ExchangeConfig struct {
	ExchangeId ExchangeId `yaml:"exchange"`
	ApiUrl     string     `yaml:"api_url"`

	ApiKeyRef     secret.Secret `yaml:"api_key"`
	SecretKeyRef  secret.Secret `yaml:"secret_key"`
	PassphraseRef secret.Secret `yaml:"passphrase"`

	// Loads all api_key, secret_key, passphrase if this is set.
	// The format must either be:
	// - separated by newlines in order of (<api_key>, <secret_key>, [passphrase])
	// - or a JSON object with the keys "api_key", "secret_key", and optionally "passphrase"
	SecretsRef secret.Secret `yaml:"secrets"`
}

func (c *ExchangeConfig) LoadSecrets() error {
	if c.SecretsRef == "" {
		return nil
	}
	slog.Debug("loading secrets for", "exchange", c.ExchangeId, "secrets", c.SecretsRef)
	value, err := c.SecretsRef.Load()
	if err != nil {
		return err
	}
	jsonMaybe := map[string]string{}
	jsonErr := json.Unmarshal([]byte(value), &jsonMaybe)
	if jsonErr != nil {
		// try splitting by newlines
		lines := strings.Split(value, "\n")
		if len(lines) < 2 {
			return fmt.Errorf("expected 2 or 3 lines of secrets, got %d", len(lines))
		}
		c.ApiKeyRef = secret.NewRawSecret(strings.TrimSpace(lines[0]))
		c.SecretKeyRef = secret.NewRawSecret(strings.TrimSpace(lines[1]))
		if len(lines) > 2 {
			c.PassphraseRef = secret.NewRawSecret(strings.TrimSpace(lines[2]))
		}
	} else {
		if apiKey, ok := jsonMaybe["api_key"]; ok {
			c.ApiKeyRef = secret.NewRawSecret(apiKey)
		} else {
			return fmt.Errorf("api_key is required")
		}
		if secretKey, ok := jsonMaybe["secret_key"]; ok {
			c.SecretKeyRef = secret.NewRawSecret(secretKey)
		} else {
			return fmt.Errorf("secret_key is required")
		}
		if passphrase, ok := jsonMaybe["passphrase"]; ok {
			c.PassphraseRef = secret.NewRawSecret(passphrase)
		}
	}
	return nil
}
