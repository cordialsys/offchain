package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/cordialsys/offchain/cmd"
	"github.com/cordialsys/offchain/cmd/oc/exchange"
	"github.com/cordialsys/offchain/cmd/oc/keys"
	"github.com/spf13/cobra"
)

func printJson(data interface{}) {
	bz, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bz))
}

func NewRootCmd() *cobra.Command {
	var verbose int
	cmd := &cobra.Command{
		SilenceUsage: true,
		PersistentPreRunE: func(cmdPre *cobra.Command, args []string) error {
			cmd.SetVerbosity(verbose)

			ctx := cmdPre.Context()
			cmdPre.SetContext(ctx)

			return nil
		},
	}
	cmd.AddCommand(NewConfigCmd())
	cmd.AddCommand(NewSecretCmd())
	cmd.AddCommand(NewStartCmd())
	cmd.AddCommand(exchange.NewExchangeCmd())
	cmd.AddCommand(exchange.NewOffchainClientCmd())
	cmd.AddCommand(keys.NewKeysCmd())
	cmd.PersistentFlags().CountVarP(
		&verbose,
		"verbose",
		"v",
		"the verbosity level",
	)

	return cmd
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
