package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	oc "github.com/cordialsys/offchain"
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

func loadFromEnv() *oc.ExchangeConfig {
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	passphrase := os.Getenv("PASSPHRASE")
	return &oc.ExchangeConfig{
		ExchangeId: oc.Okx,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		Passphrase: passphrase,
	}
}
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
	}
	cmd.AddCommand(NewGetAssetsCmd())
	cmd.AddCommand(NewGetBalancesCmd())
	cmd.AddCommand(NewAccountTransferCmd())
	cmd.AddCommand(NewWithdrawCmd())
	return cmd
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
