package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zampsn/ztool/internal/config"
	"log/slog"
)

func Config() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		RunE:  viewConfig,
	}
	cmd.AddCommand(createConfigCmd, viewConfigCmd)
	return cmd
}

var createConfigCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a configuration",
	RunE:  createConfig,
}

var viewConfigCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	RunE:  viewConfig,
}

func createConfig(cmd *cobra.Command, args []string) error {
	return config.InteractiveCreate()
}

func viewConfig(cmd *cobra.Command, args []string) error {
	cfg := viper.AllSettings()
	b, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}
	slog.Info(string(b))
	return nil
}
