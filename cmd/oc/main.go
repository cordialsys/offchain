package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func printJson(data interface{}) {
	bz, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bz))
}

const XcSolAddress = "8Rub84DA2L6BH2FVKzVBxD7jp1i6o1BtrCunm51hQcg2"
const OkxSolAddress = "CHMcJTRFxjBpcNTsfkXtckjVCqEyR6ESAhuV9tFQC3WE"

type contextKey string

const exchangeConfigKey contextKey = "exchange_config"
const configContextKey contextKey = "config"

func unwrapExchangeConfig(ctx context.Context) *oc.ExchangeConfig {
	return ctx.Value(exchangeConfigKey).(*oc.ExchangeConfig)
}

func NewRootCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config, err := loader.LoadConfig(configPath)
			if err != nil {
				return err
			}
			exchangeConfig, ok := config.GetExchange(oc.Okx)
			if !ok {
				return fmt.Errorf("exchange not found")
			}

			ctx := cmd.Context()
			ctx = context.WithValue(ctx, exchangeConfigKey, exchangeConfig)
			ctx = context.WithValue(ctx, configContextKey, config)
			cmd.SetContext(ctx)

			return nil
		},
	}
	cmd.AddCommand(NewGetAssetsCmd())
	cmd.AddCommand(NewGetBalancesCmd())
	cmd.AddCommand(NewAccountTransferCmd())
	cmd.AddCommand(NewWithdrawCmd())
	cmd.AddCommand(NewConfigCmd())

	cmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		"",
		fmt.Sprintf("path to the config file (may set %s)", loader.ENV_OFFCHAIN_CONFIG),
	)

	return cmd
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
