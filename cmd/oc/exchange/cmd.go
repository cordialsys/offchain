package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/cmd"
	"github.com/cordialsys/offchain/exchanges/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/pkg/httpsignature/signer"
	"github.com/cordialsys/offchain/pkg/keyring"
	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/cordialsys/offchain/server/client"
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

const exchangeClientKey contextKey = "exchange_client"

func unwrapClient(ctx context.Context) loader.Client {
	return ctx.Value(exchangeClientKey).(loader.Client)
}

func wrapClient(ctx context.Context, client loader.Client) context.Context {
	return context.WithValue(ctx, exchangeClientKey, client)
}

func addExchangeCommands(cmd *cobra.Command) {
	cmd.AddCommand(NewGetAssetsCmd())
	cmd.AddCommand(NewListBalancesCmd())
	cmd.AddCommand(NewAccountTransferCmd())
	cmd.AddCommand(NewWithdrawCmd())
	cmd.AddCommand(NewGetDepositAddressCmd())
	cmd.AddCommand(NewListWithdrawalHistoryCmd())
	cmd.AddCommand(NewListSubaccountsCmd())
	cmd.AddCommand(NewListAccountTypesCmd())
}

func exchangePreRun(preCmd *cobra.Command, args []string) error {
	exchange := preCmd.Flag("exchange").Value.String()
	configPath := preCmd.Flag("config").Value.String()
	subaccountId := preCmd.Flag("subaccount").Value.String()

	cmd.SetVerbosityFromCmd(preCmd)
	config, err := loader.LoadValidatedConfig(configPath)
	if err != nil {
		return err
	}

	if exchange == "" {
		return fmt.Errorf("--exchange is required")
	}
	exchangeConfig, ok := config.GetExchange(oc.ExchangeId(exchange))
	if !ok {
		return fmt.Errorf("exchange not found. options are: %v", oc.ValidExchangeIds)
	}
	slog.Info("Using exchange", "exchange", exchange)
	// ctx = context.WithValue(ctx, exchangeConfigKey, exchangeConfig)
	var account *oc.Account

	needsAuth := true
	// hack to not require auth for these commands that basically just read the config
	switch preCmd.Name() {
	case "account-types", "subaccounts":
		needsAuth = false
	}

	if needsAuth {
		if subaccountId != "" {
			for _, subaccount := range exchangeConfig.SubAccounts {
				if string(subaccount.Id) == subaccountId || subaccount.Alias == subaccountId {
					account = subaccount.AsAccount()
				}
			}
			if account == nil {
				return fmt.Errorf("subaccount %s not found in configuration for %s", subaccountId, exchange)
			}
			err = account.LoadSecrets()
			if err != nil {
				return fmt.Errorf("could not load secrets for %s subaccount %s: %w", exchange, subaccountId, err)
			}
		} else {
			account = exchangeConfig.AsAccount()
			err = account.LoadSecrets()
			if err != nil {
				return fmt.Errorf("could not load secrets for %s: %w", exchange, err)
			}
		}
	} else {
		account = NopAccount
	}
	cli, err := loader.NewClient(exchangeConfig, account)
	if err != nil {
		return fmt.Errorf("could not create client for %s: %w", exchange, err)
	}
	ctx := wrapClient(preCmd.Context(), cli)
	preCmd.SetContext(ctx)

	return nil
}

func offchainPreRun(preCmd *cobra.Command, args []string) error {
	cmd.SetVerbosityFromCmd(preCmd)
	apiUrl := preCmd.Flag("api").Value.String()
	// configPath := preCmd.Flag("config").Value.String()
	subaccountId := preCmd.Flag("subaccount").Value.String()
	signWith := preCmd.Flag("sign-with").Value.String()
	exchange := preCmd.Flag("exchange").Value.String()
	_bearerToken := preCmd.Flag("bearer-token").Value.String()
	bearerToken := secret.Secret(_bearerToken)

	keyringDir, err := preCmd.Flags().GetString("keyring-dir")
	if err != nil {
		return err
	}
	kr := keyring.New(keyring.KeyringDirOrTreasuryHome(keyringDir))

	clientOptions := []client.ClientOption{}

	if subaccountId != "" {
		clientOptions = append(clientOptions, client.WithSubAccount(subaccountId))
	}

	if signWith != "" {
		if _, err := os.Stat(signWith); err != nil {
			signWith = strings.TrimPrefix(signWith, "client-keys/")
		}
		key, err := kr.Load(signWith)
		if err != nil {
			return fmt.Errorf("could not load key %s: %w", signWith, err)
		}
		signer, err := signer.NewEd25519Signer(string(key.Secret))
		if err != nil {
			return fmt.Errorf("could not create ed25519 signer for %s: %w", signWith, err)
		}
		clientOptions = append(clientOptions, client.WithSigner(signer))
	} else {
		// try to use a bearer token
		token, err := bearerToken.Load()
		if err != nil {
			return fmt.Errorf("could not load bearer token: %w", err)
		}
		if token != "" {
			clientOptions = append(clientOptions, client.WithBearerToken(token))
		}
	}
	cli := client.NewClient(apiUrl, clientOptions...)
	offchainClient, err := offchain.NewClient(cli, oc.ExchangeId(exchange))
	if err != nil {
		return fmt.Errorf("could not create offchain client: %w", err)
	}

	ctx := preCmd.Context()
	ctx = wrapClient(ctx, offchainClient)

	preCmd.SetContext(ctx)

	return nil
}

func NewExchangeCmd() *cobra.Command {
	var configPath string
	var exchange string
	var subaccountId string
	cmd := &cobra.Command{
		Use:               "exchange",
		Short:             "Execute command directly using exchange APIs and local secrets configuration",
		Aliases:           []string{"x", "ex"},
		SilenceUsage:      true,
		PersistentPreRunE: exchangePreRun,
	}
	addExchangeCommands(cmd)

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

func NewOffchainClientCmd() *cobra.Command {
	var configPath string
	var exchange string
	var subaccountId string
	var apiUrl string
	var keyringDir string
	var signWith string
	var apiKey string
	cmd := &cobra.Command{
		Use:               "api",
		Short:             "Interact with a running offchain server",
		Aliases:           []string{"a", "offchain"},
		SilenceUsage:      true,
		PersistentPreRunE: offchainPreRun,
	}
	addExchangeCommands(cmd)

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

	cmd.PersistentFlags().StringVarP(
		&apiUrl,
		"api",
		"a",
		"http://127.0.0.1:6333",
		"The API URL to use. Defaults to the local offchain server.",
	)

	cmd.PersistentFlags().StringVarP(
		&keyringDir,
		"keyring-dir",
		"k",
		"",
		"The directory to read keyring files from.",
	)

	cmd.PersistentFlags().StringVar(
		&signWith,
		"sign-with",
		"",
		"The ID of the key (or path of key file) to sign with.",
	)

	cmd.PersistentFlags().StringVar(
		&apiKey,
		"bearer-token",
		"env:BEARER_TOKEN",
		"Set a secret reference for a bearer token for the API (only valid for read endpoints).",
	)
	return cmd
}
