package main

import (
	"fmt"
	"log"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MqttSubClient struct {
	MqttClientBase
	Streams []Stream
}

func CreateSubClient(host string, streams []Stream, clientID, username, password string) MqttSubClient {
	subClient := MqttSubClient{}
	subClient.Host = host
	subClient.Streams = streams
	subClient.Username = username
	subClient.Password = password
	subClient.Connecting = false
	subClient.Client = CreatePahoClient(host, clientID, username, password)

	return subClient
}

func (m *MqttSubClient) Start() {
	log.Printf("Starting MQTT subscription client on %s", m.Host)
	m.connect()

	if len(m.Streams) > 0 {
		for _, s := range m.Streams {
			if token := m.Client.Subscribe(s.IncomingTopic, 0, func(client paho.Client, msg paho.Message) {
				go m.handleIncomingMessage(msg.Topic(), msg.Payload(), s.OutgoingTopic)
			}); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
			}
		}
	}
}

func (m *MqttSubClient) handleIncomingMessage(topic string, payload []byte, outgoingTopic string) {
	log.Printf("Topic: %v, Msg, %v CONVERTING MESSAGE TO SENSORTHINGS", topic, string(payload[:]))
}
