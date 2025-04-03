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

type Auth struct {
	Read  []BearerToken `yaml:"read"`
	Write []BearerToken `yaml:"write"`
}
type Config struct {
	Listen    string   `yaml:"listen" env-default:":6333"`
	Origins   []string `yaml:"origins" env-default:""`
	AnyOrigin bool     `yaml:"any_origin" env-default:"false"`
	Auth      Auth     `yaml:"auth"`
}

const ENV_OFFCHAIN_CONFIG = "OFFCHAIN_CONFIG"

func LoadConfig(configPathMaybe string) (*Config, error) {
	type configsection struct {
		Offchain Config `yaml:"server"`
	}
	section := configsection{}
	if configPathMaybe == "" {
		if v := os.Getenv(ENV_OFFCHAIN_CONFIG); v != "" {
			configPathMaybe = v
		}
	}
	if configPathMaybe == "" {
		return nil, fmt.Errorf("config path is required (maybe set %s)", ENV_OFFCHAIN_CONFIG)
	}

	logrus.WithField("config", configPathMaybe).Info("loading configuration")
	err := cleanenv.ReadConfig(configPathMaybe, &section)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration: %v", err)
	}

	return &section.Offchain, nil
}
