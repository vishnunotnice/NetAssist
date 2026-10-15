package cli

import (
	"github.com/spf13/cobra"
	"github.com/vishnunotnice/NetAssist/internal/discovery"
)

func newDiscoverCommand() *cobra.Command {
	var maxHosts = 256
	cmd := &cobra.Command{
		Use:   "discover <local-ipv4-cidr>",
		Short: "Plan a bounded discovery run for a local network",
		Long:  "Lists the local addresses that an explicitly approved future probe may inspect. This command does not send network traffic.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addresses, err := discovery.Plan(args[0], maxHosts)
			if err != nil {
				return err
			}
			for _, address := range addresses {
				if err := writef(cmd.OutOrStdout(), "%s\n", address); err != nil {
					return err
				}
			}
			return writef(cmd.OutOrStdout(), "Planned %d host(s); no probes were sent.\n", len(addresses))
		},
	}
	cmd.Flags().IntVar(&maxHosts, "max-hosts", 256, "maximum eligible hosts to include")
	return cmd
}
