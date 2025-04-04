package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func SetVerbosityFromCmd(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetCount("verbose")
	if err != nil {
		panic(err)
	}
	SetVerbosity(verbose)
}
func SetVerbosity(verbose int) {
	switch verbose {
	case 0:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case 1:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case 2:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}
