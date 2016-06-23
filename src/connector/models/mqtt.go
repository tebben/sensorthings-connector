package models

import "time"

// PubClient defines the publish paho client that sends to the SubBroker
//   ClientID: the name used to connect to the publish and subscription brokers
//   Qos: quality of service, used when pubbing and subbing, defaults to 0
//   KeepAlive: will set the amount of time (in seconds) that the client
// 	should wait before sending a PING request to the broker. This will
// 	allow the client to know that a connection has not been lost with the
// 	server.
//   PingTimeout: will set the amount of time (in seconds) that the client
//	will wait after sending a PING request to the broker, before deciding
//	that the connection has been lost
type PubClient struct {
	ClientID    string        `json:"clientId"`
	Qos         byte          `json:"qos"`
	KeepAlive   time.Duration `json:"keepAlive"`
	PingTimeOut time.Duration `json:"pingTimeout"`
}

// PubBroker defines the publish broker where the incoming messages will be published to
//   Host: Host of the broker including scheme (tcp, ssl, ws), ip or hostname and port, for instance tcp://iot.eclipse.org:1883
//   Username: username needed to connect to the broker, leave blank or remove if not needed
//   Password: password needed to connect to the broker, leave blank or remove if not needed
type PubBroker struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// SubBroker defines a subscription broker
type SubBroker struct {
	ClientID string   `json:"clientId"`
	QOS      byte     `json:"qos"`
	Host     string   `json:"host"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Streams  []Stream `json:"streams"`
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
