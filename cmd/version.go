package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	passman.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("passman CLI password manager 0.0.1")
	},
}
