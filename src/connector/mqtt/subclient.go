package mqtt

import (
	"encoding/json"
	"log"
	"strconv"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"time"
)

// MqttSubClient is the implementation of the subscription client, the subscription client
// will connect to a broker where messages can be received
type MqttSubClient struct {
	MqttClientBase
	Streams []models.Stream
}

// CreateSubClient instantiates a MqttSubClient
func CreateSubClient(host string, qos byte, streams []models.Stream, clientID string, channel chan *models.PublishMessage, username string, password string, keepAlive time.Duration, pingTimeout time.Duration) MqttSubClient {
	subClient := MqttSubClient{}
	subClient.SetClientBase(host, qos, clientID, channel, username, password, keepAlive, pingTimeout)
	subClient.Streams = streams
	return subClient
}

// Start will start the subscription client by connecting and subscribing on topics
func (m *MqttSubClient) Start() {
	log.Printf("Starting MQTT subscription client on %s", m.Host)
	m.connect()

	if len(m.Streams) > 0 {
		for idx, s := range m.Streams {
			st := m.Streams[idx]
			if token := m.Client.Subscribe(s.IncomingTopic, m.Qos, func(client paho.Client, msg paho.Message) {
				go m.handleIncomingMessage(msg.Topic(), msg.Payload(), st.Mapping, st.OutgoingTopic)
			}); token.Wait() && token.Error() != nil {
				log.Print(token.Error())
			}
		}
	}
}

// handleIncomingMessage handles an incoming message by converting the payload into a message thet can be used in a
// SensorThings server and sending it over the PublishChannel to the publish client
func (m *MqttSubClient) handleIncomingMessage(topic string, payload []byte, mapping map[string]models.ToValue, outgoingTopic string) {
	if len(mapping) == 0 {
		return
	}

	var msg map[string]interface{}
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return
	}

	o := &models.Observation{}
	for k, v := range mapping {
		incVal, ok := msg[k]
		if ok {
			switch v.Name {
			case "result":
				{
					if v.ToFloat {
						iValue, err := strconv.ParseFloat(incVal.(string), 64)
						if err == nil {
							o.Result = iValue
						} else {
							o.Result = incVal
						}
					} else {
						o.Result = incVal
					}
				}
			case "phenomenonTime":
				{
					o.PhenomenonTime = incVal.(string)
				}
			}
		}
	}

	m.PublishChannel <- &models.PublishMessage{Topic: outgoingTopic, Observation: o}
}
