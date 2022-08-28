package passmanCache

import (
	"os"
	"fmt"
	"sort"
	"crypto/sha256"
	"errors"
	"bytes"
	"encoding/json"
	"passman/pkg/passmanCrypt"
	"passman/pkg/passmanConfig"
)

type Cache struct {
	Seed []byte
	Keywords []string
}

func (cache *Cache) AddKeyword(word string) {
	cache.Keywords = append(cache.Keywords, word)
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
	data,decryptErr := passmanCrypt.Decrypt(bytes.NewBuffer(cryptedData), key)
	if decryptErr != nil {
		return Cache{}, errors.New("could not decrypt cache with provided key")
	}
	var cache Cache
	if jsonErr := json.Unmarshal(data.Bytes(), &cache); jsonErr != nil {
		fmt.Println(jsonErr)
		return Cache{}, errors.New("Json is corrupted")
	}
	return cache, nil
}

func (cache *Cache) Save(key []byte) (error) {
	config,_ := passmanConfig.GetConfig()
	file, err := os.Create(config.CacheLocation)
	if err != nil {
		return errors.New("Could not open " + config.CacheLocation)
	}
	defer file.Close()
	json,err := json.Marshal(cache)
	if err != nil {
		fmt.Println(err)
		return errors.New("Could not get cache in json format")
	}
	cryptedData,cryptErr := passmanCrypt.Crypt(bytes.NewBuffer(json), key)
	if cryptErr != nil {
		fmt.Println(err)
		return errors.New("Could not crypt data")
	}
	_,err = file.Write(cryptedData.Bytes())
	if err != nil {
		file.Close()
		return errors.New("Could not write to " + config.CacheLocation)
	}
	os.Chmod(config.CacheLocation, 0600)
	return nil
}

