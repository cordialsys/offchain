package offchain

import (
	"encoding/json"
	"fmt"
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

type MultiSecret struct {
	ApiKeyRef     secret.Secret `yaml:"api_key"`
	SecretKeyRef  secret.Secret `yaml:"secret_key"`
	PassphraseRef secret.Secret `yaml:"passphrase"`

	// Loads all api_key, secret_key, passphrase if this is set.
	// The format must either be:
	// - separated by newlines in order of (<api_key>, <secret_key>, [passphrase])
	// - or a JSON object with the keys "api_key", "secret_key", and optionally "passphrase"
	SecretsRef secret.Secret `yaml:"secrets"`
}

type ExchangeClientConfig struct {
	ApiUrl string `yaml:"api_url"`
}

// Main configuration for an exchange
type ExchangeConfig struct {
	ExchangeId           ExchangeId `yaml:"exchange"`
	ExchangeClientConfig `yaml:",inline"`

	// ApiKeyRef     secret.Secret `yaml:"api_key"`
	// SecretKeyRef  secret.Secret `yaml:"secret_key"`
	// PassphraseRef secret.Secret `yaml:"passphrase"`
	MultiSecret `yaml:",inline"`

	// Subaccounts are isolated accounts on an exchange.  They have their own API keys and are
	// typically used for trading.
	SubAccounts []*SubAccount `yaml:"subaccounts"`
}

// A subaccount is an isolated account on an exchange.  It's just like a main account, but typically you
// cannot withdraw funds from it directly.  Instead, funds must be transferred to the main account first.
// Subaccounts have their own independent API keys.
type SubAccount struct {
	// The id of the subaccount is required.  This is created on the exchange by the user.
	// If it does not match, then it won't work.
	Id          string `yaml:"id"`
	MultiSecret `yaml:",inline"`
}

func (c *MultiSecret) LoadSecrets() error {
	if c.SecretsRef == "" {
		return nil
	}
	value, err := c.SecretsRef.Load()
	if err != nil {
		return err
	}
	jsonMaybe := map[string]string{}
	jsonErr := json.Unmarshal([]byte(value), &jsonMaybe)
	if jsonErr != nil {
		// try splitting by newlines
		lines := strings.Split(value, "\n")

		lines = func() []string {
			trimmed := []string{}
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				trimmed = append(trimmed, line)
			}
			return trimmed
		}()
		if len(lines) == 0 {
			return fmt.Errorf("secret value is empty")
		}

		// If the first line is in "<api-key>:<secret-key>", then convert it to be two separate lines
		lines = func() []string {
			first := lines[0]
			if strings.Count(first, ":") == 1 {
				parts := strings.Split(first, ":")
				newParts := []string{parts[0], parts[1]}
				newParts = append(newParts, lines[1:]...)
				return newParts
			}
			return lines
		}()
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
