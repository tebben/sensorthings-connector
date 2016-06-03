package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	PubBroker  PubBroker   `json:"pubBroker"`
	SubBrokers []SubBroker `json:"subBrokers"`
}

type PubBroker struct {
	Info BrokerInfo `json:"info"`
}

type SubBroker struct {
	Info    BrokerInfo `json:"info"`
	Streams []Stream   `json:"streams"`
}

type BrokerInfo struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Stream struct {
	IncomingTopic string `json:"topicIn"`
	OutgoingTopic string `json:"topicOut"`
}

type Convert struct {
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
