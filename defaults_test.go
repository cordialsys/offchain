package offchain_test

import (
	"slices"
	"testing"

	oc "github.com/cordialsys/offchain"
	"github.com/stretchr/testify/require"
)

var validAliases = []string{
	"spot",
	"trading",
	"funding",
	"derivatives",
	"cross-margin",
	"isolated-margin",
}

func TestDefaults(t *testing.T) {
	for _, exchange := range oc.ValidExchangeIds {
		cfg, ok := oc.GetDefaultConfig(exchange)
		if !ok {
			require.Failf(t, "no account types found for exchange", "exchange=%s", exchange)
		}
		at := cfg.AccountTypes
		if cfg.NoAccountTypes != nil && *cfg.NoAccountTypes {
			require.Empty(t, at, "account types defined for exchange but not expected", "exchange=%s", exchange)
			continue
		} else {
			require.NotEmpty(t, at, "no account types defined for exchange %s", exchange)
		}

		for _, at := range at {
			require.NotEmpty(t, at.Type, "account type must be defined for exchange %s", exchange)
		}

		aliases := map[string]bool{}
		for _, at := range at {
			for _, alias := range at.Aliases {
				if _, ok := aliases[alias]; ok {
					require.Fail(t, "aliases must be unique for exchange %s, but %s is duplicated", exchange, alias)
				}
				aliases[alias] = true

				if !slices.Contains(validAliases, alias) {
					require.Fail(t,
						"new defined alias for exchange %s, alias=%s.  Consider using an existing alias before adding it as a valid one.",
						exchange, alias,
					)
				}
			}
		}

		accountTypes := map[string]bool{}
		for _, at := range at {
			accountTypes[string(at.Type)] = true
		}
		require.Equal(t, len(accountTypes), len(at), "account types must be unique for exchange=%s", exchange)

	}
}
