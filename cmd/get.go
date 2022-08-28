package cmd

import (
	"fmt"
	"syscall"
	"os"
	"errors"
	"passman/pkg/passmanCache"
	"passman/pkg/passmanCrypt"
	"golang.org/x/term"
	"github.com/spf13/cobra"
	"github.com/atotto/clipboard"
)

const getDescription = `Bwill give you the password generated from the given key. By default the key 
will be put into your clipboard but use -p to use stdout`

func init() {
	passman.AddCommand(getCmd)
}

var getCmd = &cobra.Command {
	Use: "get",
	Short: "will give you the password generated from the given key",
	Long: getDescription,
	Run: getExec,
}

func getExec(cmd *cobra.Command, args []string) {
	fmt.Printf("password: ")
	password,_ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	cache,cacheErr := passmanCache.GetCache(password)
	if cacheErr != nil {
		fmt.Println("Password is not valid")
		os.Exit(1)
	}
	genPass := passmanCrypt.GenPass(cache.Seed, []byte(args[0]))
	if clipErr := useClipboard(genPass); clipErr != nil {
		os.Exit(1)
	}
	//fmt.Println(genPass)
}

func useClipboard(password string) (error){
	if err := clipboard.WriteAll(password); err != nil {
		fmt.Println(err)
		return errors.New("Could not use clipboard")
	}
	return nil
}
