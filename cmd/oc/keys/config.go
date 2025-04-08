package keys

import (
	"fmt"

	"github.com/cordialsys/offchain/pkg/keyring"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show the keyring location",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyringDir, err := cmd.Flags().GetString("keyring-dir")
			if err != nil {
				return err
			}
			fmt.Println(keyring.KeyringDirOrTreasuryHome(keyringDir))
			return nil
		},
	}

	return cmd
}
