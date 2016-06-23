package mqtt

import (
	"encoding/json"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"log"
)

// MqttPubClient is the implementation of the publish client, the publish client
// will connect to the broker where all the incoming messages needs to be send to
type MqttPubClient struct {
	MqttClientBase
}

// CreatePubClient instantiates a MqttPubClient
func CreatePubClient(host string, qos byte, clientID string, channel chan *models.PublishMessage, username, password string) MqttPubClient {
	pubClient := MqttPubClient{}
	pubClient.SetClientBase(host, qos, clientID, channel, username, password)
	return pubClient
}

// Start will start the publish client by connecting and start listening on the PublishChannel
func (m *MqttPubClient) Start() {
	log.Printf("Starting MQTT publish client on %s", m.Host)
	m.connect()
	go m.listen()
}

// listen start listening for publish messages on the PublishChannel
func (m *MqttPubClient) listen() {
	for {
		pm := <-m.PublishChannel
		if !m.Connecting {
			jsonString, err := json.Marshal(pm.Observation)
			if err != nil {
				log.Printf("Error marshalling observation: %v", err.Error())
			}

			token := m.Client.Publish(pm.Topic, m.Qos, false, jsonString)
			token.Wait()

		}
	}
}
