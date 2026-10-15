package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show local NetAssist foundation status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := os.Stat(configPath)
			configStatus := "not found"
			if err == nil {
				configStatus = "available"
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("inspect configuration: %w", err)
			}
			return writef(cmd.OutOrStdout(), "CLI: ready\nConfiguration: %s (%s)\nServices: not configured\n", configStatus, configPath)
		},
	}
	cmd.Flags().StringVar(&configPath, "config", defaultConfigPath(), "configuration file path")
	return cmd
}
