package keys

import (
	"encoding/hex"
	"fmt"

	"github.com/cordialsys/offchain/pkg/keyring"
	"github.com/spf13/cobra"
)

func NewGenerateCmd() *cobra.Command {
	var overwrite bool
	cmd := &cobra.Command{
		Use:   "generate [id]",
		Short: "Generate and store a new key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyringDir, err := cmd.Flags().GetString("keyring-dir")
			if err != nil {
				return err
			}
			kr := keyring.New(keyring.KeyringDirOrTreasuryHome(keyringDir))

			if _, err = kr.Load(args[0]); err == nil {
				if !overwrite {
					return fmt.Errorf("key already exists")
				}
			}

			k, err := kr.New(keyring.NewClientKeyName(args[0]), keyring.Ed25519)
			if err != nil {
				return err
			}
			fmt.Println(k.Name)
			pub, err := k.PublicKey()
			if err != nil {
				return err
			}
			fmt.Println(hex.EncodeToString(pub))
			return nil
		},
	}

	cmd.Flags().BoolVarP(
		&overwrite,
		"overwrite",
		"o",
		false,
		"Overwrite any existing key",
	)
	return cmd
}
