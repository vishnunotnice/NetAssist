package cli

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/vishnunotnice/NetAssist/internal/devices"
)

func newDevicesCommand() *cobra.Command {
	var storePath string
	cmd := &cobra.Command{Use: "devices", Short: "Manage the local device inventory"}
	cmd.PersistentFlags().StringVar(&storePath, "store", defaultDataPath("devices.json"), "device inventory path")
	cmd.AddCommand(newDeviceAddCommand(&storePath), newDeviceListCommand(&storePath))
	return cmd
}

func newDeviceAddCommand(storePath *string) *cobra.Command {
	var address, tags string
	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a managed device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			device, err := devices.New(args[0], address, strings.Split(tags, ","))
			if err != nil {
				return err
			}
			if err := (devices.FileStore{Path: *storePath}).Add(device); err != nil {
				return err
			}
			return writef(cmd.OutOrStdout(), "Added %s (%s)\n", device.Name, device.ID)
		},
	}
	cmd.Flags().StringVar(&address, "address", "", "device hostname or IP address")
	cmd.Flags().StringVar(&tags, "tags", "", "comma-separated device tags")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}

func newDeviceListCommand(storePath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List managed devices",
		RunE: func(cmd *cobra.Command, _ []string) error {
			items, err := (devices.FileStore{Path: *storePath}).List()
			if err != nil {
				return err
			}
			for _, device := range items {
				if err := writef(cmd.OutOrStdout(), "%s\t%s\t%s\n", device.Name, device.Address, strings.Join(device.Tags, ",")); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
