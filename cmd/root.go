package cmd

import (
	"os"

	"github.com/analog-substance/krypton/internal/run"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krypton",
	Short: "An arsenic style tool to scan internal networks via embedded scripts",
	RunE: func(cmd *cobra.Command, args []string) error {
		networks, _ := cmd.Flags().GetString("networks")
		fromDisk, _ := cmd.Flags().GetBool("from-disk")

		if fromDisk {
			run.SetExecMode(run.ExecFromDisk)
		}

		hosts, err := run.DiscoverHosts(networks)
		if err != nil {
			return err
		}

		err = run.DiscoverTCPServices(hosts)
		if err != nil {
			return err
		}

		return run.DiscoverUDPServices(hosts)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("networks", "n", "", "The CIDR networks on which to discover hosts. Defaults to the machines current networks.")
	rootCmd.Flags().Bool("from-disk", false, "Executes the scripts from disk instead of from memory")
}
