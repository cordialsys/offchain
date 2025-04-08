package keyring

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const TREASURY_HOME = "TREASURY_HOME"

type ClientKeyName string

const Prefix = "client-keys"

func NewClientKeyName(id string) ClientKeyName {
	id = strings.TrimSuffix(id, ".toml")
	return ClientKeyName(filepath.Join(Prefix, id))
}
func (n ClientKeyName) Id() string {
	parts := strings.Split(string(n), "/")
	return parts[len(parts)-1]
}

type Algorithm string

const (
	EcdsaK256Sha256 Algorithm = "ecdsa-k256-sha256"
	Ed25519         Algorithm = "ed25519"
)

type Hex string

func (h Hex) Decode() ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(string(h), "0x"))
}

type ClientKey struct {
	Name      ClientKeyName `toml:"name"`
	Algorithm Algorithm     `toml:"algorithm"`
	Secret    Hex           `toml:"secret_key"`
}

// Returns the public key
func (ck *ClientKey) PublicKey() ([]byte, error) {
	switch ck.Algorithm {
	case Ed25519:
		secretBz, err := ck.Secret.Decode()
		if err != nil {
			return nil, err
		}
		keypair := ed25519.NewKeyFromSeed(secretBz)
		return keypair.Public().(ed25519.PublicKey), nil
	default:
		return nil, fmt.Errorf("unsupported alg: %s", ck.Algorithm)
	}
}

type Keyring struct {
	dir string
}

const DefaultSubDir = "keyring"

// Create a new keyring, passing the directory
// that contains the .toml key files, typically $XDG_DATA_HOME/treasury/$TREASURY_ID or $TREASURY_HOME
func New(dir string) Keyring {
	return Keyring{dir}
}

func (keyring *Keyring) Load(filename string) (*ClientKey, error) {
	ck := &ClientKey{}
	bz, err := os.ReadFile(filepath.Join(keyring.dir, filename))
	if errors.Is(err, os.ErrNotExist) {
		// try with .toml
		bz, err = os.ReadFile(filepath.Join(keyring.dir, filename+".toml"))
	}
	if err != nil {
		return ck, err
	}
	err = toml.Unmarshal(bz, ck)
	if err != nil {
		return ck, err
	}
	return ck, nil
}

func (keyring *Keyring) List() ([]string, error) {
	files, err := os.ReadDir(keyring.dir)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".toml") {
			id := strings.TrimSuffix(file.Name(), ".toml")
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (keyring *Keyring) New(name ClientKeyName, alg Algorithm) (*ClientKey, error) {
	var secret []byte
	switch alg {
	case Ed25519, "ed255":
		secret = GenerateEd255Seed()
	default:
		return nil, fmt.Errorf("unsupported alg: %s", alg)
	}
	return keyring.Import(name, alg, secret)
}

func (keyring *Keyring) Import(name ClientKeyName, alg Algorithm, secret []byte) (*ClientKey, error) {
	ck := &ClientKey{
		Name:      name,
		Algorithm: alg,
		Secret:    Hex(hex.EncodeToString(secret)),
	}
	bz, err := toml.Marshal(ck)
	if err != nil {
		return ck, err
	}
	// fmt.Println("opening", keyring.dir, "for", name)
	err = os.MkdirAll(keyring.dir, os.ModePerm)
	if err != nil {
		return ck, err
	}
	// fmt.Println("writing", filepath.Join(keyring.dir, name.Id()+".toml"))
	err = os.WriteFile(
		filepath.Join(keyring.dir, name.Id()+".toml"),
		bz,
		os.ModePerm,
	)
	if err != nil {
		return ck, err
	}
	return ck, nil
}

func GenerateEd255Seed() []byte {
	bz := make([]byte, 32)
	_, err := rand.Read(bz)
	if err != nil {
		panic(err)
	}
	return bz
}

// Return the keyring from the input path.
// If there is no directory in the path (id or filename only),
// Then the default keyring path based on TreasuryHome is used.
func KeyringDirOrTreasuryHome(path string) string {
	dirMaybe := filepath.Dir(path)
	dirMaybe = replaceTilda(dirMaybe)

	// 1. TREASURY_HOME
	// 2. XDG_CONFIG_HOME/offchain/$TREASURY_ID
	// 3. $HOME/.offchain/$TREASURY_ID
	if dirMaybe == "." {
		home := os.Getenv(TREASURY_HOME)
		if home == "" {
			home, _ = os.UserConfigDir()
			home = filepath.Join(home, "offchain")
			if home == "" {
				home, _ = os.UserHomeDir()
				home = filepath.Join(home, ".offchain")
			}
		}
		dirMaybe = filepath.Join(home, DefaultSubDir)
	}
	return dirMaybe
}
func replaceTilda(path string) string {
	if len(path) > 1 && path[0] == '~' {
		userHome, _ := os.UserHomeDir()
		path = strings.Replace(path, "~", userHome, 1)
	}
	return path
}
