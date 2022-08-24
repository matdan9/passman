package passmanCache

import (
	"os"
	"fmt"
	"sort"
	"crypto/sha256"
	"errors"
	"encoding/json"
	"passman/pkg/passmanCrypt"
	"passman/pkg/passmanConfig"
)

type Cache struct {
	Seed []byte
	keywords []string
}

func (cache *Cache) GenerateSeed(words []string) {
	sort.Strings(words)
	var seed string
	for i:=0; i<len(words); i++ { 
		seed += words[i]
	}
	hashedSeed := sha256.Sum256([]byte(seed))
	cache.Seed = hashedSeed[:]
}

func GetCache(key []byte) (Cache, error) {
	config,_ := passmanConfig.GetConfig()
	file,err := os.Open(config.CacheLocation)
	if err != nil {
		return Cache{}, errors.New("Looks like passman is not setup yet, use \"passman init\" to get started")
	}
	fInfo,_ := file.Stat()
	cryptedData := make([]byte, fInfo.Size())
	if _,err := file.Read(cryptedData); err != nil {
		return Cache{}, errors.New("Could not read file " + config.CacheLocation)
	}
	if err = file.Close(); err != nil {
		return Cache{}, errors.New("Could not close " + config.CacheLocation + " properly")
	}
	data,decryptErr := passmanCrypt.Decrypt(cryptedData, key)
	if decryptErr != nil {
		fmt.Println(decryptErr);
		return Cache{}, errors.New("could not decrypt cache with provided key")
	}
	var cache Cache
	if jsonErr := json.Unmarshal(data, &cache); jsonErr != nil {
		fmt.Println(jsonErr)
		errors.New("Json if corrupted")
	}
	return cache, nil
}

func (cache *Cache) Save(key []byte) (error) {
	config,_ := passmanConfig.GetConfig()
	file, err := os.Create(config.CacheLocation)
	if err != nil {
		return errors.New("Could not open " + config.CacheLocation)
	}
	json,err := json.Marshal(cache)
	if err != nil {
		fmt.Println(err)
		return errors.New("Could not get cache in json format")
	}
	cryptedData := passmanCrypt.Crypt(json, key)
	_,err = file.Write(cryptedData)
	if err != nil {
		file.Close()
		return errors.New("Could not write to " + config.CacheLocation)
	}
	if err = file.Close(); err != nil {
		return errors.New("Could not close cache file")
	}
	os.Chmod(config.CacheLocation, 0600)
	return nil
}

