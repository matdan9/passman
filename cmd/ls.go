package cmd

import (
	"fmt"
	"syscall"
	"os"
	"passman/pkg/passmanCache"
	"golang.org/x/term"
	"github.com/spf13/cobra"
)

const lsDescription= `Will list the previosly saved keywords`

func init() {
	passman.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command {
	Use: "ls",
	Short: "Will list saved keywords",
	Long: lsDescription,
	Run: lsExec,
}

func lsExec(cmd *cobra.Command, args []string) {
	fmt.Printf("password: ")
	password,_ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println("\n")
	cache,cacheErr := passmanCache.GetCache(password)
	if cacheErr != nil {
		fmt.Println("Password is not valid")
		os.Exit(1)
	}
	for _,e := range cache.Keywords {
		fmt.Println(e)
	}
	fmt.Printf("\nTotal: %d\n", len(cache.Keywords))
}
