package main

import (
	"fmt"
	"log/slog"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/pkg/httpsignature/verifier"
	"github.com/cordialsys/offchain/server"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	var listen string
	var configPath string
	var anyOrigin bool
	var origins []string
	var publicReadEndpoints bool
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Serve the offchain server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loader.LoadValidatedConfig(configPath)
			if err != nil {
				return err
			}
			serverConfig, err := server.LoadConfig(configPath)
			if err != nil {
				return err
			}
			if listen != "" {
				serverConfig.Listen = listen
			}
			if anyOrigin {
				serverConfig.AnyOrigin = true
			}
			if publicReadEndpoints {
				serverConfig.PublicReadEndpoints = true
			}
			if len(origins) > 0 {
				serverConfig.Origins = origins
			}
			bearers := make([]string, len(serverConfig.BearerTokens))
			for i, bearer := range serverConfig.BearerTokens {
				bearers[i], err = bearer.Token.Load()
				if err != nil {
					return fmt.Errorf("failed to load bearer token %s: %w", bearer.Id, err)
				}
			}
			publicKeys := make([]verifier.VerifierI, len(serverConfig.PublicKeys))
			for i, key := range serverConfig.PublicKeys {
				switch key.Algorithm {
				case server.Ed25519, "":
					publicKeys[i], err = verifier.NewEd25519Verifier(key.Key.Bytes())
					if err != nil {
						return fmt.Errorf("failed to create ed25519 verifier: %w", err)
					}
				default:
					return fmt.Errorf("unsupported algorithm: %s", key.Algorithm)
				}
			}

			if len(bearers) == 0 && len(publicKeys) == 0 && !publicReadEndpoints {
				slog.Warn("no authentication means configured, read-only endpoints will be unreachable")
			}
			if len(publicKeys) == 0 {
				slog.Warn("no public keys configured, write-endpoints will be unreachable")
			}

			serverArgs := server.ServerArgs{
				Listen:              serverConfig.Listen,
				AnyOrigin:           serverConfig.AnyOrigin,
				Origins:             serverConfig.Origins,
				BearerTokens:        bearers,
				PublicKeys:          publicKeys,
				PublicReadEndpoints: serverConfig.PublicReadEndpoints,
			}
			server := server.New(config, serverArgs)
			return server.Start()
		},
	}
	cmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		fmt.Sprintf("path to the config file (may set %s)", oc.ENV_OFFCHAIN_CONFIG),
	)
	cmd.Flags().StringVarP(&listen, "listen", "l", "", "address to listen on (default is 127.0.0.1:6333)")
	cmd.Flags().BoolVar(&anyOrigin, "any-origin", false, "allow any origin for CORS")
	cmd.Flags().StringSliceVar(&origins, "origins", []string{}, "origins to allow for CORS")
	cmd.Flags().BoolVar(&publicReadEndpoints, "public-read-endpoints", false, "Disable authentication for read endpoints")
	return cmd
}
