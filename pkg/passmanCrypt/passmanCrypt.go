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

func Crypt(inBuff *bytes.Buffer, key []byte) (*bytes.Buffer, error) {
	cryptedBuff := bytes.NewBuffer([]byte{})
	// TODO predict the final length of the cryptedBuff and Grow it to the accoring size
	block := getBlock(key)
	iv := make([]byte, block.BlockSize())
	if _,err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not crypt data")
	}
	cryptedBuff.Write(iv)
	stream := cipher.NewCTR(block, iv)
	err := processFullBuffer(inBuff, cryptedBuff, &stream)
	return cryptedBuff, err
}

func Decrypt(inBuff *bytes.Buffer, key []byte) (*bytes.Buffer, error) {
	outBuff := bytes.NewBuffer([]byte{})
	block := getBlock(key)
	iv := make([]byte, block.BlockSize())
	if _,err := inBuff.Read(iv); err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not read provided buffer")
	}
	stream := cipher.NewCTR(block, iv)
	err := processFullBuffer(inBuff, outBuff, &stream)
	return outBuff, err
}

func processFullBuffer(inBuff *bytes.Buffer, outBuff *bytes.Buffer, stream *cipher.Stream) (error) {
	tmpBuff := make([]byte, 2048)
	for {
		n,err := inBuff.Read(tmpBuff)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return errors.New("Could not crypt data")
		}
		if n>0 {
			(*stream).XORKeyStream(tmpBuff, tmpBuff[:n])
			outBuff.Write(tmpBuff[:n])
		}
	}
	return nil
}

func getBlock (key []byte) (cipher.Block){
	hashKey := sha256.Sum256(key)
	c, err := aes.NewCipher(hashKey[:])
	if err != nil {
		fmt.Println(err)
	}
	return c
}

func GenPass(seed []byte, keyword []byte) (string) {
	dump := append(seed, keyword...)
	pass := sha256.Sum256(dump)
	return base64.StdEncoding.EncodeToString(pass[:])
}

func generateIV() ([]byte, error) {
	config,_ := passmanConfig.GetConfig()
	salt := make([]byte, config.SaltLength)
	_,err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return salt, nil
}
