package cmd

import (
	"fmt"
	"os"
	"runtime"

	static "github.com/analog-substance/arsenic-static"
	"github.com/analog-substance/krypton/internal/bin"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "krypton",
	Short: "An arsenic style tool to scan internal networks via embedded scripts",
	RunE: func(cmd *cobra.Command, args []string) error {
		networks, _ := cmd.Flags().GetString("networks")

		nmapPath, err := bin.Locate("nmap")
		if err != nil {
			nmapPath = "./nmap"

			fmt.Println("[-] nmap either not found or an error occurred. Falling back writing to disk")

			err = bin.WriteAs(fmt.Sprintf("nmap_%s", runtime.GOARCH), "nmap")
			if err != nil {
				return fmt.Errorf("error occurred while writing nmap to disk: %v", err)
			}
		}

		discoverCmd, err := static.Command("bin/as-recon-discover-hosts", networks)
		if err != nil {
			return err
		}
		discoverCmd.Env = append(discoverCmd.Env, fmt.Sprintf("NMAP=%s", nmapPath))

		hosts, err := discoverCmd.Output()
		if err != nil {
			return fmt.Errorf("error occurred while discovering hosts: %v", err)
		}

		servicesCmd, err := static.Command("bin/as-recon-discover-services")
		if err != nil {
			return err
		}
		servicesCmd.Stdout = os.Stdout

		env := []string{
			fmt.Sprintf("SCRIPT_STDIN=%s", string(hosts)),
			fmt.Sprintf("NMAP=%s", nmapPath),
		}
		servicesCmd.Env = append(servicesCmd.Env, env...)

		return servicesCmd.Run()
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
}
