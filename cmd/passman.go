package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var passman = &cobra.Command{
  Use:   "passman",
  Short: "CLI password manager",
  Long: `work in progess`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := passman.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
