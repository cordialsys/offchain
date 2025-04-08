package keys

import (
	"github.com/spf13/cobra"
)

func NewKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "keys",
		Short:        "Manage keys for authentication to offchain server",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().String(
		"keyring-dir",
		"",
		"The directory to manage keyring files",
	)

	cmd.AddCommand(NewGenerateCmd())
	cmd.AddCommand(NewConfigCmd())
	cmd.AddCommand(NewLsCmd())
	cmd.AddCommand(NewGetCmd())

	return cmd
}
