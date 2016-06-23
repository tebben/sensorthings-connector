package config

import (
	"encoding/json"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"io/ioutil"
)

// Config defines all the settings to setup the connector
//   HttpHost: the host were the rest interface should run on
//   PubClient: te publish client, see PubClient
//   PubBroker: te publish broker, see PubBroker
type Config struct {
	HttpHost  string           `json:"httpHost"`
	PubClient models.PubClient `json:"publishClient"`
	PubBroker models.PubBroker `json:"publishBroker"`
}

// readFile reads the bytes from a given file
func readFile(cfgFile string) ([]byte, error) {
	source, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	return source, nil
}

// readConfig tries to parse the byte data into a config file
func readConfig(fileContent []byte) (Config, error) {
	config := Config{}
	err := json.Unmarshal(fileContent, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// GetConfig retrieves a new configuration from the given config file
// Fatal when config does not exist or cannot be read
func GetConfig(cfgFile string) (Config, error) {
	content, err := readFile(cfgFile)
	if err != nil {
		return Config{}, err
	}

	conf, err := readConfig(content)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
