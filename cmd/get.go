package cmd

import (
	"fmt"
	"syscall"

	"passman/pkg/passmanCache"
	"passman/pkg/passmanCrypt"
	"golang.org/x/term"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

const long = `Bwill give you the password generated from the given key. By default the key 
will be put into your clipboard but use -p to use stdout`

func init() {
	passman.AddCommand(getCmd)
}

var getCmd = &cobra.Command {
	Use: "get",
	Short: "will give you the password generated from the given key",
	Long: long,
	Run: getExec,
}

func getExec(cmd *cobra.Command, args []string) {
	fmt.Printf("password: ")
	password,_ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	cache,cacheErr := passmanCache.GetCache(password)
	if cacheErr != nil {
		fmt.Println("Password is not valid")
		return
	}
	genPass := passmanCrypt.GenPass(cache.Seed, []byte(args[0]))
	useClipboard(genPass)
	fmt.Println(genPass)
}

func useClipboard(password string) {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Could not use clipboard")
	}
	clipboard.Write(clipboard.FmtText, []byte(password))
}
