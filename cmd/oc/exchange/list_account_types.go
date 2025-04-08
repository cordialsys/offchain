package exchange

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/spf13/cobra"
)

var NopAccount = &oc.Account{
	Id:         "",
	SubAccount: false,
	MultiSecret: oc.MultiSecret{
		ApiKeyRef:     secret.Secret("raw:AAA"),
		SecretKeyRef:  secret.Secret("raw:AAA"),
		PassphraseRef: secret.Secret("raw:AAA"),
	},
}

// This just reads the configuration currently, not yet exposing an API call for this.
func NewListAccountTypesCmd() *cobra.Command {
	return &cobra.Command{
		SilenceUsage: true,
		Use:          "account-types",
		Short:        "List valid account types for an exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := unwrapClient(cmd.Context())
			accountTypes, err := cli.ListAccountTypes()
			if err != nil {
				return err
			}
			printJson(accountTypes)
			return nil
		},
	}
}
