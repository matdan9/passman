package cmd

import (
	"fmt"
	"syscall"
	"os"
	"passman/pkg/passmanCache"
	"golang.org/x/term"
	"github.com/spf13/cobra"
)

const saveDescription= `Saves the given keyword to be listed later with the ls command`

func init() {
	passman.AddCommand(saveCmd)
}

var saveCmd = &cobra.Command {
	Use: "save",
	Short: "Will save the given keyword",
	Long: saveDescription,
	Run: saveExec,
}

func saveExec(cmd *cobra.Command, args []string) {
	if(len(args) < 1) {
		os.Exit(1)
	}
	fmt.Printf("password: ")
	password,_ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	cache,cacheErr := passmanCache.GetCache(password)
	if cacheErr != nil {
		fmt.Println("Password is not valid")
		os.Exit(1)
	}
	cache.AddKeyword(args[0])
	if saveErr := cache.Save(password); saveErr != nil {
		fmt.Println("Could not save keyword")
		os.Exit(1)
	}
}
