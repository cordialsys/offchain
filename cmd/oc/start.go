package main

import (
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	var listen string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Serve the offchain server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cmd.Context().Value(configContextKey).(*loader.Config)
			server := server.New(listen, config)
			return server.Start()
		},
	}

	cmd.Flags().StringVarP(&listen, "listen", "l", ":8080", "address to listen on")

	return cmd
}
