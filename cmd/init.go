package cmd

import (
	"fmt"
	"bufio"
	"os"
	"syscall"
	"strings"
	"bytes"

	"golang.org/x/term"

	"passman/pkg/passmanCache"
	"passman/pkg/passmanConfig"

	"github.com/spf13/cobra"
)

const seedIntroduction = `We will ask you for an unlimited amount of words that will define your seed.
a seed is what allows the correct saving and creation of your passwords so for 2 computers
to share the same passwords, both of them must have the same seed aka the words you are about
to give use. After each word you must press the RETURN key and type "!DONE" when you are done.`

const doneTrigger = "!DONE"

func init() {
	passman.AddCommand(initCmd)
}

var initCmd = &cobra.Command {
	Use: "init",
	Short: "Initializes passman",
	Long: "Starts an interactive CLI to setup your passman \n" + seedIntroduction,
	Run: initExec,
}

func initExec(cmd *cobra.Command, args []string) {
	fmt.Println(seedIntroduction + "\n")
	words := readWords()
	var cache passmanCache.Cache
	cache.GenerateSeed(words);
	if saveErr := cache.Save(readConfirmPassword()); saveErr != nil {
		fmt.Println(saveErr)
	}
}

func readConfirmPassword() ([]byte) {
	var password  []byte
	var confPassword []byte 
	for true {
		fmt.Printf("Password: ")
		password,_ = term.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\nConfrim Password: ")
		confPassword,_ = term.ReadPassword(int(syscall.Stdin))
		if eq := bytes.Compare(password, confPassword); eq == 0 {
			fmt.Println("\n")
			return password
		}
		fmt.Printf("\nThe passwords are not equal\n")
	}
	return nil
}

func readWords() ([] string) {
	config,_ := passmanConfig.GetConfig()
	reader := bufio.NewReader(os.Stdin)
	var word string
	var words []string
	for i:=1; i < config.MaxWordCount; i++ {
		fmt.Printf("%d: ", i)
		word,_ = reader.ReadString('\n')
		word = strings.Replace(word, "\n", "", -1)
		if word == doneTrigger {
			if i > config.MinWordCount {
				return words
			}
			i--
			fmt.Printf("You need at least %d words\n", config.MinWordCount)
		}
		words = append(words, word)
		fmt.Println()
	}
	fmt.Printf("Reached maximum word count of %d\n", config.MaxWordCount)
	return words;
}
