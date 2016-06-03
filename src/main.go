package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

const ClientID = "sensorthings-connector"

var publishBroker MqttPubClient
var subscribeBrokers []MqttSubClient

func main() {
	cfgFlag := flag.String("config", "sampleconfig.json", "path of the config file")
	flag.Parse()
	cfg := *cfgFlag

	c, err := GetConfig(cfg)
	if err != nil {
		log.Fatal("config read error: ", err)
		return
	}

	createBrokerClients(c)

	// keep alive method for now
	t := time.NewTicker(15 * time.Minute)
	for now := range t.C {
		fmt.Sprintf("tick: %v", now)
	}
}

func createBrokerClients(c Config) {
	publishBroker = CreatePubClient(c.PubBroker.Info.Host, ClientID, c.PubBroker.Info.Username, c.PubBroker.Info.Password)
	publishBroker.Start()

	for _, sb := range c.SubBrokers {
		subCient := CreateSubClient(sb.Info.Host, sb.Streams, ClientID, sb.Info.Username, sb.Info.Password)
		subCient.Start()

		subscribeBrokers = append(subscribeBrokers, subCient)
	}
}
