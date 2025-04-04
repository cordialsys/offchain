package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/cmd"
	"github.com/spf13/cobra"
)

func printJson(data interface{}) {
	bz, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bz))
}

type contextKey string

const exchangeConfigKey contextKey = "exchange_config"
const exchangeAccountSecretsKey contextKey = "exchange_account_secrets"
const configContextKey contextKey = "config"

func unwrapExchangeConfig(ctx context.Context) *oc.ExchangeConfig {
	return ctx.Value(exchangeConfigKey).(*oc.ExchangeConfig)
}

func unwrapAccountSecrets(ctx context.Context) *oc.MultiSecret {
	return ctx.Value(exchangeAccountSecretsKey).(*oc.MultiSecret)
}

func unwrapAccountConfig(ctx context.Context) (*oc.ExchangeConfig, *oc.MultiSecret) {
	exchangeConfig := unwrapExchangeConfig(ctx)
	secrets := unwrapAccountSecrets(ctx)
	return exchangeConfig, secrets
}

func NewExchangeCmd() *cobra.Command {
	var configPath string
	var exchange string
	var subaccountId string
	cmd := &cobra.Command{
		Use:          "exchange",
		Short:        "Run commands on an exchange",
		Aliases:      []string{"ex", "x", "exchanges"},
		SilenceUsage: true,
		PersistentPreRunE: func(preCmd *cobra.Command, args []string) error {
			cmd.SetVerbosityFromCmd(preCmd)
			config, err := oc.LoadConfig(configPath)
			if err != nil {
				return err
			}
			ctx := preCmd.Context()

			if exchange == "" {
				return fmt.Errorf("--exchange is required")
			}
			exchangeConfig, ok := config.GetExchange(oc.ExchangeId(exchange))
			if !ok {
				return fmt.Errorf("exchange not found. options are: %v", oc.ValidExchangeIds)
			}
			slog.Info("Using exchange", "exchange", exchange)
			ctx = context.WithValue(ctx, exchangeConfigKey, exchangeConfig)
			var secrets *oc.MultiSecret

			if subaccountId != "" {
				for _, subaccount := range exchangeConfig.SubAccounts {
					if subaccount.Id == subaccountId {
						secrets = &subaccount.MultiSecret
					}
				}
				if secrets == nil {
					return fmt.Errorf("subaccount %s not found in configuration for %s", subaccountId, exchange)
				}
				err = secrets.LoadSecrets()
				if err != nil {
					return fmt.Errorf("could not load secrets for %s subaccount %s: %w", exchange, subaccountId, err)
				}
			} else {
				secrets = &exchangeConfig.MultiSecret
				err = secrets.LoadSecrets()
				if err != nil {
					return fmt.Errorf("could not load secrets for %s: %w", exchange, err)
				}
			}

			ctx = context.WithValue(ctx, configContextKey, config)
			ctx = context.WithValue(ctx, exchangeAccountSecretsKey, secrets)
			preCmd.SetContext(ctx)

			return nil
		},
	}
	cmd.AddCommand(NewGetAssetsCmd())
	cmd.AddCommand(NewListBalancesCmd())
	cmd.AddCommand(NewAccountTransferCmd())
	cmd.AddCommand(NewWithdrawCmd())
	cmd.AddCommand(NewGetDepositAddressCmd())
	cmd.AddCommand(NewListWithdrawalHistoryCmd())
	cmd.AddCommand(NewListSubaccountsCmd())
	cmd.AddCommand(NewListAccountTypesCmd())
	cmd.PersistentFlags().StringVarP(
		&exchange,
		"exchange",
		"x",
		"",
		fmt.Sprintf("The exchange to use (%v)", oc.ValidExchangeIds),
	)

	cmd.PersistentFlags().StringVarP(
		&subaccountId,
		"subaccount",
		"s",
		"",
		"The subaccount to use. Defaults to using the main account.",
	)

	cmd.PersistentFlags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		fmt.Sprintf("path to the config file (may set %s)", oc.ENV_OFFCHAIN_CONFIG),
	)

	return cmd
}
