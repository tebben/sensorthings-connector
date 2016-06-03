package main

import (
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient interface {
	Start()
	Stop()
	connect()
	retryConnect()
}

type MqttClientBase struct {
	Client     paho.Client
	Host       string
	Username   string
	Password   string
	Connecting bool
}

func CreatePahoClient(host string, clientID string, username string, password string) paho.Client {
	opts := paho.NewClientOptions().AddBroker(host).SetClientID(clientID)
	opts.SetKeepAlive(300 * time.Second)
	opts.SetPingTimeout(20 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(connectHandler)

	if len(username) > 0 && len(password) > 0 {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	return paho.NewClient(opts)
}

// Stop the MQTT client
func (m *MqttClientBase) Stop() {
	m.Client.Disconnect(500)
}

func (m *MqttClientBase) connect() {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		if !m.Connecting {
			log.Printf("MQTT client %v", token.Error())
			m.retryConnect()
		}
	}
}

// retryConnect starts a ticker which tries to connect every xx seconds and stops the ticker
// when a connection is established. This is useful when MQTT Broker and GOST are hosted on the same
// machine and GOST is started before mosquito
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

func connectHandler(c paho.Client) {
	log.Printf("MQTT client connected")
}

func connectionLostHandler(c paho.Client, err error) {
	log.Printf("MQTT client lost connection: %v", err)
}
