package main

import (
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// MqttClient defines the needed methods to control our MQTT client
type MqttClient interface {
	Start()
	Stop()
	connect()
	retryConnect()
}

// MqttClientBase holds information on our client needed for a
// Publish and Subscription client
type MqttClientBase struct {
	Qos            byte
	Client         paho.Client
	Host           string
	Username       string
	Password       string
	Connecting     bool
	PublishChannel chan *PublishMessage
}

// SetClientBase sets the base parameters needed for our MQTT client and creates the MQTT client
func (m *MqttClientBase) SetClientBase(host string, qos byte, clientID string, channel chan *PublishMessage, username, password string) {
	m.Qos = qos
	m.Host = host
	m.Username = username
	m.Password = password
	m.Connecting = false
	m.Client = createPahoClient(host, clientID, username, password)
	m.PublishChannel = channel
}

// Stop will stop the MQTT client
func (m *MqttClientBase) Stop() {
	m.Client.Disconnect(500)
}

// connect tries to connect the MQTT client to the broker, on fail the retry procedure kicks in
func (m *MqttClientBase) connect() {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		if !m.Connecting {
			log.Printf("MQTT client %v", token.Error())
			m.retryConnect()
		}
	}
}

// retryConnect starts a ticker which tries to connect every xx seconds and stops the ticker
// when a connection is established.
func (m *MqttClientBase) retryConnect() {
	log.Printf("MQTT client %s starting reconnect procedure in background", m.Host)

	m.Connecting = true
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			m.connect()
			if m.Client.IsConnected() {
				ticker.Stop()
				m.Connecting = false
			}
		}
	}()
}

// createPahoClient creates a new paho client
func createPahoClient(host string, clientID string, username string, password string) paho.Client {
	opts := paho.NewClientOptions().AddBroker(host).SetClientID(clientID)
	opts.SetKeepAlive(300 * time.Second)
	opts.SetPingTimeout(20 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(func(client paho.Client, err error) {
		connectionLostHandler(client, err, host)
	})
	opts.SetOnConnectHandler(func(client paho.Client) {
		connectHandler(client, host)
	})

	if len(username) > 0 && len(password) > 0 {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	return paho.NewClient(opts)
}

func connectHandler(c paho.Client, host string) {
	log.Printf("MQTT client connected on %s", host)
}

func connectionLostHandler(c paho.Client, err error, host string) {
	log.Printf("MQTT client lost connection on: %s: errorL%v", host, err)
}
