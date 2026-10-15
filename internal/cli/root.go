// Package cli defines the public NetAssist command-line interface.
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const applicationName = "netassist"

// BuildVersion is set at build time with -ldflags "-X .../cli.BuildVersion=vX.Y.Z".
var BuildVersion = "dev"

// NewRootCommand returns an isolated root command, suitable for production and tests.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          applicationName,
		Short:        "Network operations and device management",
		Long:         "NetAssist is a foundation for network management, diagnostics, monitoring, and security event operations.",
		SilenceUsage: true,
	}
	cmd.AddCommand(newVersionCommand(), newStatusCommand(), newConfigCommand(), newDiagnosticsCommand(), newDiscoverCommand(), newDevicesCommand(), newMonitorCommand(), newSIEMCommand())
	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writef(cmd.OutOrStdout(), "%s %s\n", applicationName, BuildVersion)
		},
	}
}

func writef(w io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(w, format, args...)
	return err
}
