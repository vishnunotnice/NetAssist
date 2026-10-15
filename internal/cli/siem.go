package cli

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/vishnunotnice/NetAssist/internal/siem"
)

func newSIEMCommand() *cobra.Command {
	var storePath string
	cmd := &cobra.Command{Use: "siem", Short: "Manage local security events"}
	cmd.PersistentFlags().StringVar(&storePath, "store", defaultDataPath("siem-events.jsonl"), "SIEM event store path")
	cmd.AddCommand(newSIEMIngestCommand(&storePath), newSIEMEventsCommand(&storePath))
	return cmd
}

func newSIEMIngestCommand(storePath *string) *cobra.Command {
	var source, kind, severity, message string
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Record one normalized security event",
		RunE: func(cmd *cobra.Command, _ []string) error {
			event, err := siem.NewEvent(source, kind, severity, message)
			if err != nil {
				return err
			}
			if err := (siem.FileStore{Path: *storePath}).Append(event); err != nil {
				return err
			}
			return writef(cmd.OutOrStdout(), "Recorded %s\n", event.ID)
		},
	}
	cmd.Flags().StringVar(&source, "source", "", "event source")
	cmd.Flags().StringVar(&kind, "kind", "", "event kind")
	cmd.Flags().StringVar(&severity, "severity", "info", "info, low, medium, high, or critical")
	cmd.Flags().StringVar(&message, "message", "", "event message")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("kind")
	_ = cmd.MarkFlagRequired("message")
	return cmd
}

func newSIEMEventsCommand(storePath *string) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "events",
		Short: "List locally stored security events",
		RunE: func(cmd *cobra.Command, _ []string) error {
			events, err := (siem.FileStore{Path: *storePath}).List(limit)
			if err != nil {
				return err
			}
			for _, event := range events {
				if err := writef(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\t%s\n", event.Timestamp.Format(time.RFC3339), event.Severity, event.Source, event.Kind, event.Message); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 20, "maximum events to display; zero shows all")
	return cmd
}
