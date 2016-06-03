package main

import "log"

type MqttPubClient struct {
	MqttClientBase
}

func CreatePubClient(host string, clientID, username, password string) MqttPubClient {
	pubClient := MqttPubClient{}
	pubClient.Host = host
	pubClient.Username = username
	pubClient.Password = password
	pubClient.Connecting = false
	pubClient.Client = CreatePahoClient(host, clientID, username, password)

	return pubClient
}

func (m *MqttPubClient) Start() {
	log.Printf("Starting MQTT subscription client on %s", m.Host)
	m.connect()
}
