package offchain

import "github.com/cordialsys/offchain/pkg/secret"

type ExchangeConfig struct {
	ExchangeId ExchangeId `yaml:"exchange"`
	ApiUrl     string     `yaml:"api_url"`

	ApiKeyRef     secret.Secret `yaml:"api_key" env-default:"env:API_KEY"`
	SecretKeyRef  secret.Secret `yaml:"secret_key" env-default:"env:SECRET_KEY"`
	PassphraseRef secret.Secret `yaml:"passphrase" env-default:"env:PASSPHRASE"`
}
