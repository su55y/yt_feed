package config

import (
	"io/ioutil"
	"log"
	"regexp"
	"sync"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	API_KEY    string `yaml:"api_key"`
	ApiKeyPath string `yaml:"api_key_path"`
	// alternative cache path, overrides default if directory exists
	CachePath  string   `yaml:"cache_dir"`
	MaxResults int64    `yaml:"max_results"`
	Region     string   `yaml:"region"`
	ThumbOff   bool     `yaml:"thumbnails_disable"`
	ThumbSize  string   `yaml:"thumbnails_size"`
	Channels   []string `yaml:"channels"`
	ThumbDir   string
}

var (
	confInstance     AppConfig
	once             sync.Once
	unmarshalError   error
	channelIdPattern = regexp.MustCompile("^[a-zA-Z0-9\\-_]{24}$")
)

func GetAppConfig(path string) (AppConfig, error) {
	once.Do(func() {
		confInstance = AppConfig{}
		unmarshalError = yaml.Unmarshal(readFile(path), &confInstance)
		if unmarshalError != nil {
			log.Printf("config unmarshal error: %v\n", unmarshalError)
		}

		// filter invalid channel ids
		ids := make([]string, 0)
		for _, i := range confInstance.Channels {
			if channelIdPattern.MatchString(i) {
				ids = append(ids, i)
			}
		}
		confInstance.Channels = ids
	})

	return confInstance, unmarshalError
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("config load error %v\n", err)
	}
	return data
}
