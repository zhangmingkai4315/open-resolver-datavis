package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VERSION define the version number
const VERSION = "v0.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of open-resolver-datavis",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Open-Resolver-Datavis Version: %s\n", VERSION)
	},
}
