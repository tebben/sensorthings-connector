package mqtt

import (
	"encoding/json"
	"errors"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	connectorMQTT "github.com/tebben/sensorthings-connector/src/connector/mqtt"
)

// MQTTModule can be used as middleware to connect a MQTT stream of observations (using structured data) to a
// SensorThings server.
type MQTTModule struct {
	models.ConnectorModuleBase
	settings   MQTTModuleSettings
	subClients []connectorMQTT.MqttSubClient
}

// MQTTModuleSettings is used to configure the listening MQTT clients
//   SubBrokers can contain information on multiple MQTT brokers where the module should subscribe to
type MQTTModuleSettings struct {
	SubBrokers []models.SubBroker `json:"subBrokers"`
}

// Setup initialised the module by setting some default values
func (mq *MQTTModule) Setup() {
	mq.Name = "MQTT"
	mq.Description = "Map a structured non MQTT stream to a SensorThings observation stream"
}

// Start will create MQTT subscription clients that are configured in the settings and start them
func (mq *MQTTModule) Start() {
	mq.subClients = []connectorMQTT.MqttSubClient{}
	for _, sb := range mq.settings.SubBrokers {
		subClient := connectorMQTT.CreateSubClient(sb.Host, sb.QOS, sb.Streams, sb.ClientID, mq.PublishChannel, sb.Username, sb.Password, 300, 20)
		subClient.Start()

		mq.subClients = append(mq.subClients, subClient)
	}
}

// Stop will stop all defined Subscription brokers
func (mq *MQTTModule) Stop() {
	if mq.subClients == nil {
		return
	}

	for idx := range mq.subClients {
		mq.subClients[idx].Stop()
	}
}

// SettingsChanged will try to parse and set MQTTModuleSettings from a json.RawMessage
func (mq *MQTTModule) SettingsChanged(settings json.RawMessage) error {
	s := MQTTModuleSettings{}
	if err := json.Unmarshal(settings, &s); err != nil {
		return errors.New("Unable to read MQTT Module settings")
	}

	mq.settings = s
	return nil
}
