package cli

import (
	"context"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/vishnunotnice/NetAssist/internal/diagnostics"
)

func newDiagnosticsCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "diagnose", Short: "Run explicit-target network diagnostics"}
	cmd.AddCommand(newDNSCommand(), newTCPCommand())
	return cmd
}

func newDNSCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dns <hostname>",
		Short: "Resolve a hostname",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := diagnostics.LookupDNS(cmd.Context(), net.DefaultResolver, args[0])
			if err != nil {
				return err
			}
			for _, address := range result.Addresses {
				if err := writef(cmd.OutOrStdout(), "%s\n", address); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newTCPCommand() *cobra.Command {
	var timeout time.Duration
	cmd := &cobra.Command{
		Use:   "tcp <host:port>",
		Short: "Check one TCP endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()
			result, err := diagnostics.CheckTCP(ctx, args[0], timeout)
			if err != nil {
				return err
			}
			state := "unreachable"
			if result.Reachable {
				state = "reachable"
			}
			return writef(cmd.OutOrStdout(), "%s: %s (%s)\n", result.Address, state, result.Duration.Round(time.Millisecond))
		},
	}
	cmd.Flags().DurationVar(&timeout, "timeout", 3*time.Second, "connection timeout")
	return cmd
}
