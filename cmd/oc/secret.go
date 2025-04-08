package main

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/spf13/cobra"
)

func NewSecretCmd() *cobra.Command {
	multi := false
	help := ""
	for _, t := range secret.Types {
		help += fmt.Sprintf("%s: %s\n", t, t.Usage())
	}
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "secret <reference>",
		Short:        "Print out a secret given a reference",
		Long:         help,
		Args:         cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if multi {
				multiSecret := oc.MultiSecret{
					SecretsRef: secret.Secret(args[0]),
				}
				err = multiSecret.LoadSecrets()
				if err != nil {
					return err
				}
				printJson(multiSecret)
			} else {
				sec, err := secret.Secret(args[0]).Load()
				if err != nil {
					return err
				}
				fmt.Println(sec)
			}

			return nil
		},
	}
	cmd.Flags().BoolVar(&multi, "multi", false, "Parse the secret value, expecting multi-part exchange secrets")
	return cmd
}
