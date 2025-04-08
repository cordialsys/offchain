package keys

import (
	"encoding/hex"
	"fmt"

	"github.com/cordialsys/offchain/pkg/keyring"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyringDir, err := cmd.Flags().GetString("keyring-dir")
			if err != nil {
				return err
			}
			kr := keyring.New(keyring.KeyringDirOrTreasuryHome(keyringDir))
			k, err := kr.Load(args[0])
			if err != nil {
				return err
			}
			pub, err := k.PublicKey()
			if err != nil {
				return err
			}
			fmt.Println(hex.EncodeToString(pub))
			return nil
		},
	}

	return cmd
}
