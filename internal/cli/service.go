package cli

import (
	"encoding/json"
	"fmt"

	"github.com/hi-donwi/Donwi-Brain/internal/service"
	"github.com/spf13/cobra"
)

func newServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the background lokal serve service",
		Long: `Manage the background lokal serve daemon as a systemd user service on Linux
or launchd agent on macOS. The service ensures the gateway starts on boot.`,
	}

	cmd.AddCommand(newServiceInstallCmd())
	cmd.AddCommand(newServiceUninstallCmd())
	cmd.AddCommand(newServiceStatusCmd())

	return cmd
}

func newServiceInstallCmd() *cobra.Command {
	var cloudOnly, localOnly bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install and start the background gateway service",
		Long: `Install and start the background gateway service. The service runs "serve";
--cloud-only or --local-only is passed on to it, as with "serve" itself.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cloudOnly && localOnly {
				return fmt.Errorf("cannot combine --cloud-only and --local-only")
			}
			args := []string{"serve"}
			if cloudOnly {
				args = append(args, "--cloud-only")
			}
			if localOnly {
				args = append(args, "--local-only")
			}
			res, err := service.Install(service.Options{Args: args})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "[service] %s\n", res.Message)
			return nil
		},
	}
	cmd.Flags().BoolVar(&cloudOnly, "cloud-only", false, "Run the service with serve --cloud-only (provider accounts, no local inference backend)")
	cmd.Flags().BoolVar(&localOnly, "local-only", false, "Run the service with serve --local-only (configured local backend, no cloud accounts)")
	return cmd
}

func newServiceUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Stop and remove the background gateway service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := service.Uninstall(service.Options{})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "[service] %s\n", res.Message)
			return nil
		},
	}
}

func newServiceStatusCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the status of the background gateway service",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			stat, err := service.Status(service.Options{})
			if err != nil {
				return err
			}
			if asJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(stat)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "[service] %s\n", stat.Message)
			fmt.Fprintf(cmd.OutOrStdout(), "  Type:      %s\n", stat.ServiceType)
			fmt.Fprintf(cmd.OutOrStdout(), "  Installed: %t (%s)\n", stat.Installed, stat.Path)
			fmt.Fprintf(cmd.OutOrStdout(), "  Running:   %t\n", stat.Running)
			if stat.Running && stat.Version != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "  Version:   %s\n", stat.Version)
				fmt.Fprintf(cmd.OutOrStdout(), "  Profile:   %s\n", stat.Profile)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "output status in JSON format")
	return cmd
}
