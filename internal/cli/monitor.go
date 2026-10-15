package cli

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/vishnunotnice/NetAssist/internal/monitoring"
)

func newMonitorCommand() *cobra.Command {
	var historyPath string
	cmd := &cobra.Command{Use: "monitor", Short: "Run and review explicit health checks"}
	cmd.PersistentFlags().StringVar(&historyPath, "history", defaultDataPath("monitoring.jsonl"), "monitoring history path")
	cmd.AddCommand(newMonitorTCPCommand(&historyPath), newMonitorHistoryCommand(&historyPath))
	return cmd
}

func newMonitorTCPCommand(historyPath *string) *cobra.Command {
	var timeout time.Duration
	cmd := &cobra.Command{
		Use:   "tcp <host:port>",
		Short: "Check and record one TCP endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()
			check, err := monitoring.CheckTCP(ctx, args[0], timeout, monitoring.FileHistory{Path: *historyPath})
			if err != nil {
				return err
			}
			return writef(cmd.OutOrStdout(), "%s\t%t\t%s\n", check.Target, check.Reachable, check.Duration.Round(time.Millisecond))
		},
	}
	cmd.Flags().DurationVar(&timeout, "timeout", 3*time.Second, "connection timeout")
	return cmd
}

func newMonitorHistoryCommand(historyPath *string) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "history",
		Short: "List locally recorded health checks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks, err := (monitoring.FileHistory{Path: *historyPath}).List(limit)
			if err != nil {
				return err
			}
			for _, check := range checks {
				if err := writef(cmd.OutOrStdout(), "%s\t%s\t%t\t%s\n", check.Timestamp.Format(time.RFC3339), check.Target, check.Reachable, check.Duration.Round(time.Millisecond)); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum checks to display; zero shows all")
	return cmd
}
