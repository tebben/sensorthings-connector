package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config defines all the settings to setup the connector
//   ClientID: the name used to connect to the publish and subscription brokers
//   Qos: quality of service, used when pubbing and subbing, defaults to 0
//   KeepAlive: will set the amount of time (in seconds) that the client
// 	should wait before sending a PING request to the broker. This will
// 	allow the client to know that a connection has not been lost with the
// 	server.
//   PingTimeout: will set the amount of time (in seconds) that the client
//	will wait after sending a PING request to the broker, before deciding
//	that the connection has been lost
//   PubBroker: te publish broker, see PubBroker
//   SubBrokers: the subscription brokers, possible to connect to multiple brokers see SubBroker
type Config struct {
	ClientID    string      `json:"clientId"`
	Qos         byte        `json:"qos"`
	KeepAlive   int         `json:"keepAlive"`
	PingTimeOut int         `json:"pingTimeout"`
	PubBroker   PubBroker   `json:"pubBroker"`
	SubBrokers  []SubBroker `json:"subBrokers"`
}

// PubBroker defines the publish broker where the incoming messages will be published to
type PubBroker struct {
	Info BrokerInfo `json:"info"`
}

// SubBroker defines a subscription broker
type SubBroker struct {
	Info    BrokerInfo `json:"info"`
	Streams []Stream   `json:"streams"`
}

// BrokerInfo defines the info needed to connect to a broker
//   Host: Host of the broker including scheme (tcp, ssl, ws), ip or hostname and port, for instance tcp://iot.eclipse.org:1883
//   Username: username needed to connect to the broker, leave blank or remove if not needed
//   Password: password needed to connect to the broker, leave blank or remove if not needed
type BrokerInfo struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Stream defines a datastream coming from a subscription broker
//   IncomingTopic: The topic where the connector will subscribe to
//   OutgoingTopic: The topic where the connector will publish the message to
//   Mapping: FromValue -> ToValue, simple implementation for our use-case now
type Stream struct {
	IncomingTopic string             `json:"topicIn"`
	OutgoingTopic string             `json:"topicOut"`
	Mapping       map[string]ToValue `json:"mapping"`
}

// ToValue defines the SensorThings output value, used in combination with an
// in parameter.
//   Name: the name of the SensorThings observation parameter i.e. result, phenomenonTime
//   ToInt: if the value needs to be converted into an int, useful for string values from the incoming data
type ToValue struct {
	Name    string `json:"name"`
	ToFloat bool   `json:"toFloat"`
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
