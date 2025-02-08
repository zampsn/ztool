package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	WSUtil WSUtil `json:"wsutil"`
}

type WSUtil struct {
	Addr    string            `json:"server_addr"`
	Headers map[string]string `json:"server_headers"`
}

func Init(path string) error {
	slog.Debug("Initializing config")

	// set default config options
	viper.SetDefault("wsutil.server_addr", "ws://localhost:8080")
	viper.SetDefault("wsutil.server_headers", map[string]string{})

	if path == "" {
		slog.Debug("No config file path specified, using default path")
		return loadDefaultConfig()
	} else {
		slog.Debug("Loading custom config file", "path", path)
		viper.SetConfigFile(path)
		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
	}
	return nil
}

func InteractiveCreate() error {
	path, err := getPathInput()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	slog.Debug("Creating directory", "dir", dir)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	if err = viper.SafeWriteConfigAs(path); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	// the joys of structured logging... I should use a different log package
	slog.Info(fmt.Sprintf("Config file created at: %s", path))

	// pretty json
	b, _ := json.Marshal(viper.AllSettings())
	slog.Info(string(b))
	return nil
}

func getPathInput() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	defaultPath := filepath.Join(configDir, "config.json")

	prompt := promptui.Prompt{
		Label:   "Path of the config file",
		Default: defaultPath,
	}

	path, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to prompt for path: %w", err)
	}
	return path, nil
}

func loadDefaultConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(configDir, "config.json")
	viper.SetConfigFile(configFile)
	if err = viper.MergeInConfig(); err != nil {
		var configErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configErr) {
			slog.Debug("Config file not found, using default values")
			return nil
		}
		return fmt.Errorf("could not merge load file: %w", err)
	}
	viper.Set("config", configFile)
	return nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".ztool"), nil
}
