package keyring_test

import (
	"os"
	"testing"

	"github.com/cordialsys/offchain/pkg/keyring"
	"github.com/stretchr/testify/require"
)

func TestKeyring(t *testing.T) {
	n := keyring.NewClientKeyName("id")
	require.EqualValues(t, "client-keys/id", n)
	require.Equal(t, "id", n.Id())

	n = keyring.NewClientKeyName("id.toml")
	require.EqualValues(t, "client-keys/id", n)
	require.Equal(t, "id", n.Id())

	require.Equal(t, "id", keyring.ClientKeyName("id").Id())
}

func TestGetKeyringSecret(t *testing.T) {
	os.Setenv("TREASURY_HOME", "abc")
	defer os.Unsetenv("TREASURY_HOME")
	dir := keyring.KeyringDirOrTreasuryHome("")
	require.Equal(t, dir, "abc/keyring")

	dir = keyring.KeyringDirOrTreasuryHome("id")
	require.Equal(t, dir, "abc/keyring")

	dir = keyring.KeyringDirOrTreasuryHome("otherdir/id")
	require.Equal(t, dir, "otherdir")
}
