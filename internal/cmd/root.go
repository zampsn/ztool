package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zampsn/ztool/internal/config"
	"log/slog"
	"os"
)

var rootCmd = &cobra.Command{
	Short:             "CLI Swiss Army Knife",
	Long:              "A multipurpose CLI helper for boring, repetitive, or tedious tasks",
	PersistentPreRunE: prerun,
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "path to custom configuration file)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PreRunE = prerun

	// register commands
	rootCmd.AddCommand(Config())
}

func prerun(cmd *cobra.Command, args []string) error {
	if viper.GetBool("debug") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		// Remove timestamps and excess clutter from output
		initSimpleLogger()
	}

	// Load config
	path, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to parse config flag: %w", err)
	}
	return config.Init(path)
}

func Execute() {
	_ = rootCmd.Execute()
}

type simpleLogger struct{ slog.Handler }

func (h *simpleLogger) Handle(ctx context.Context, r slog.Record) error {
	_, err := os.Stdout.WriteString(r.Message + "\n")
	return err
}

func initSimpleLogger() {
	logger := slog.New(&simpleLogger{
		Handler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		}),
	})
	slog.SetDefault(logger)
}
