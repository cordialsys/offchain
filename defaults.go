package offchain

import (
	"bytes"
	_ "embed"
	"io"

	"github.com/ilyakaznacheev/cleanenv"
)

//go:embed defaults.yaml
var defaultsYAML string

func init() {
	type configsection struct {
		Offchain Config `yaml:"offchain"`
	}
	section := configsection{}
	cleanenv.ParseYAML(io.Reader(bytes.NewBufferString(defaultsYAML)), &section)
	err := section.Offchain.Init()
	if err != nil {
		panic(err)
	}
	for _, exchage := range ValidExchangeIds {
		defaultConfigs[exchage] = &ExchangeClientConfig{}
	}
	// set account types for each exchange
	for _, exchange := range section.Offchain.Exchanges {
		// set account types for each exchange
		defaultConfigs[exchange.ExchangeId] = &exchange.ExchangeClientConfig
	}
}

var defaultConfigs = map[ExchangeId]*ExchangeClientConfig{}

func ApplyDefaults(cfg *ExchangeConfig) {
	def := defaultConfigs[cfg.ExchangeId]
	// copy over account types if defined
	if len(cfg.AccountTypes) == 0 {
		cfg.AccountTypes = make([]*AccountTypeConfig, len(def.AccountTypes))
		for i, at := range def.AccountTypes {
			cpy := *at
			cfg.AccountTypes[i] = &cpy
		}
	}

	if cfg.NoAccountTypes == nil {
		if def.NoAccountTypes != nil {
			cpy := *def.NoAccountTypes
			cfg.NoAccountTypes = &cpy
		}
	}
}

func GetDefaultConfig(exchangeId ExchangeId) (ExchangeClientConfig, bool) {
	cfg, ok := defaultConfigs[exchangeId]
	if !ok {
		return ExchangeClientConfig{}, false
	}
	return *cfg, true
}
