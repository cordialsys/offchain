package server

import (
	"fmt"
	"os"

	"github.com/cordialsys/offchain/pkg/hex"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type BearerToken struct {
	Token string `yaml:"token"`
	Name  string `yaml:"name"`
}

type Algorithm string

const (
	K256Sha256 Algorithm = "k256-sha2"
	Ed25519    Algorithm = "ed25519"
	Ed255      Algorithm = "ed255" // alias to ed25519
)

type HttpPublicKey struct {
	Algorithm Algorithm `yaml:"algorithm"`
	// Hex encoded
	Key hex.Hex `yaml:"key"`
}

type Config struct {
	Listen    string   `yaml:"listen" env-default:":6333"`
	Origins   []string `yaml:"origins" env-default:""`
	AnyOrigin bool     `yaml:"any_origin" env-default:"false"`

	BearerTokens []BearerToken   `yaml:"bearer_tokens"`
	PublicKeys   []HttpPublicKey `yaml:"public_keys"`
}

const ENV_OFFCHAIN_CONFIG = "OFFCHAIN_CONFIG"

func LoadConfig(configPathMaybe string) (*Config, error) {
	type configsection2 struct {
		Server Config `yaml:"server"`
	}
	type configsection1 struct {
		Offchain configsection2 `yaml:"offchain"`
	}

	section := configsection1{}
	if configPathMaybe == "" {
		if v := os.Getenv(ENV_OFFCHAIN_CONFIG); v != "" {
			configPathMaybe = v
		}
	}
	if configPathMaybe == "" {
		return nil, fmt.Errorf("config path is required (maybe set %s)", ENV_OFFCHAIN_CONFIG)
	}

	logrus.WithField("config", configPathMaybe).Info("loading server configuration")
	err := cleanenv.ReadConfig(configPathMaybe, &section)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration: %v", err)
	}

	return &section.Offchain.Server, nil
}
