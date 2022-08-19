package main

import (
	"passman/cmd"
	/*
	"passman/pkg/passmanCrypt"
	"fmt"
	*/
)

func main(){
	/*
	crypt := passmanCrypt.Crypt([]byte("this is a test with multiple words"), []byte("secret"))
	fmt.Printf("crypted: %s\n", crypt)
	decrypt,_ := passmanCrypt.DeCrypt([]byte(crypt), []byte("secret"))
	fmt.Printf("decrypted: %s\n", decrypt)
	wrongDecrypt,_ := passmanCrypt.DeCrypt([]byte(crypt), []byte("wrongSecret"))
	fmt.Printf("wrong decrypt: %s\n", wrongDecrypt)
	*/
	cmd.Execute()
}
