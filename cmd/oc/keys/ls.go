package keys

import (
	"fmt"

	"github.com/cordialsys/offchain/pkg/keyring"
	"github.com/spf13/cobra"
)

func NewLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			keyringDir, err := cmd.Flags().GetString("keyring-dir")
			if err != nil {
				return err
			}
			kr := keyring.New(keyring.KeyringDirOrTreasuryHome(keyringDir))
			keys, err := kr.List()
			if err != nil {
				return err
			}
			for _, key := range keys {
				fmt.Println(key)
			}
			return nil
		},
	}
	return cmd
}
