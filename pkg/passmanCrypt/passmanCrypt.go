package passmanCrypt

import (
	"io"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"encoding/base64"
	"crypto/rand"
	"bytes"
	"passman/pkg/passmanConfig"
)

func PrintPass() {
	fmt.Printf("Password Print\n")
}


func Crypt(data []byte, key []byte) ([]byte) {
	salt,_ := generateSalt()
	gcm := getAesGcm(append(key, []byte(salt)...))
	nonce := make([]byte, gcm.NonceSize())
	cryptedSeed := gcm.Seal(nonce, nonce, data, nil)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(cryptedSeed)))
	base64.StdEncoding.Encode(dst, cryptedSeed)
	dst = append(dst, []byte("\n")...)
	dst = append(dst, salt...)
	return dst
}

func DeCrypt(data []byte, key []byte) (string, error) {
	splits := bytes.Split(data, []byte("\n"))
	data = splits[0]
	salt := splits[1]
	gcm := getAesGcm(append(key, salt...))
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		fmt.Println("Wrong length")
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	base64.StdEncoding.Decode(dst, data)
	dst = dst[:len(dst)-1]
	nonce, data := dst[:nonceSize], dst[nonceSize:]
	plainData, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		// TODO change for logger output
		fmt.Println(err)
		return "", errors.New("Could not decrypt with provided key")
	}
	return string(plainData), nil
}

func getAesGcm (key []byte) (cipher.AEAD){
	hashKey := sha256.Sum256(key)
	c, err := aes.NewCipher(hashKey[:])
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	return gcm
}

func generateSalt() ([]byte, error) {
	config,_ := passmanConfig.GetConfig()
	salt := make([]byte, config.SaltLength)
	_,err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(salt)))
	base64.StdEncoding.Encode(dst, salt)
	return dst, nil
}
