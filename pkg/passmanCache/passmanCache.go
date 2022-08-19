package passmanCache

import (
	"os"
	"fmt"
	"sort"
	"crypto/sha256"
	"errors"
	"passman/pkg/passmanCrypt"
	"passman/pkg/passmanConfig"
)

func SetSeed(words []string, key []byte) (error) {
	sort.Strings(words);
	var seed string
	for i:=0; i<len(words); i++ { 
		seed += words[i]
	}
	seedBuff := sha256.Sum256([]byte(seed))
	if err := saveCache(seedBuff[:], key); err != nil {
		return err
	}
	fmt.Printf("%v\n", words);
	return nil
}

func GetCacheContent() ([]byte, error) {
	config,_ := passmanConfig.GetConfig()
	file,err := os.Open(config.CacheLocation)
	if err != nil {
		return nil, errors.New("Looks like passman is not setup yet, use \"passman init\" to get started")
	}
	fInfo,_ := file.Stat()
	data := make([]byte, fInfo.Size())
	if _,err := file.Read(data); err != nil {
		return nil, errors.New("Could not read file " + config.CacheLocation)
	}
	if err = file.Close(); err != nil {
		return nil, errors.New("Could not close " + config.CacheLocation + " properly")
	}
	return data, nil
}

func saveCache(seed []byte, key []byte) (error) {
	cryptedSeed := passmanCrypt.Crypt(seed, key)
	config,_ := passmanConfig.GetConfig()
	file, err := os.Create(config.CacheLocation)
	if err != nil {
		return errors.New("Could not open " + config.CacheLocation)
	}
	_,err = file.Write(cryptedSeed)
	if err != nil {
		file.Close()
		return errors.New("Could not write to " + config.CacheLocation)
	}
	file.Close()
	os.Chmod(config.CacheLocation, 0600)
	return nil
}

