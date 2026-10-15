package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Config contains intentionally small, stable settings for the initial CLI.
// Each future service owns its detailed runtime configuration.
type Config struct {
	Log  LogConfig  `yaml:"log"`
	SIEM SIEMConfig `yaml:"siem"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// SIEMConfig reserves configuration ownership for a future event pipeline.
// It does not enable event collection or detection behavior in Phase 1.
type SIEMConfig struct {
	Enabled bool `yaml:"enabled"`
}

func defaultConfigPath() string {
	if path, err := os.UserConfigDir(); err == nil {
		return filepath.Join(path, applicationName, "config.yaml")
	}
	return "config.yaml"
}

func defaultConfig() Config {
	return Config{Log: LogConfig{Level: "info", Format: "json"}, SIEM: SIEMConfig{Enabled: false}}
}

func newConfigCommand() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{Use: "config", Short: "Inspect or initialize local configuration"}
	cmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath(), "configuration file path")
	cmd.AddCommand(newConfigShowCommand(&configPath), newConfigInitCommand(&configPath))
	return cmd
}

func newConfigShowCommand(path *string) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the active configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			config, err := loadConfig(*path)
			if err != nil {
				return err
			}
			out, err := yaml.Marshal(config)
			if err != nil {
				return fmt.Errorf("encode configuration: %w", err)
			}
			return writef(cmd.OutOrStdout(), "%s", out)
		},
	}
}

func newConfigInitCommand(path *string) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a default configuration file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := os.MkdirAll(filepath.Dir(*path), 0o755); err != nil {
				return fmt.Errorf("create configuration directory: %w", err)
			}
			if _, err := os.Stat(*path); err == nil {
				return fmt.Errorf("configuration already exists: %s", *path)
			} else if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("inspect configuration: %w", err)
			}
			out, err := yaml.Marshal(defaultConfig())
			if err != nil {
				return fmt.Errorf("encode configuration: %w", err)
			}
			if err := os.WriteFile(*path, out, 0o600); err != nil {
				return fmt.Errorf("write configuration: %w", err)
			}
			return writef(cmd.OutOrStdout(), "Created %s\n", *path)
		},
	}
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read configuration: %w", err)
	}
	config := defaultConfig()
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("parse configuration: %w", err)
	}
	return config, nil
}
