package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"machine-operator/cmd/start"
	"machine-operator/cmd/stop"
)

var rootCmd = &cobra.Command{
	Use:          "machine-operator",
	Short:        "machine-operator a tool for actions workflow",
	SilenceUsage: true,
	Long:         `machine-operator`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(start.StartCmd)
	rootCmd.AddCommand(stop.StartCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
