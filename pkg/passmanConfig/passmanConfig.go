package passmanConfig;

import (
	"os"
)

const configFile = "passman.conf"
var config *Config

type Config struct {
	CacheLocation string
	SaltLength int
	MinWordCount int
	MaxWordCount int
}

func GetConfig() (Config, error) {
	if config != nil {
		return *config, nil
	}
	//TODO read config from config file
	config := &Config {
		MinWordCount: 2,
		MaxWordCount: 25,
		SaltLength: 50,
	}
	dir,_ := os.UserCacheDir()
	config.CacheLocation = dir + "/passman"
	return *config, nil
}

