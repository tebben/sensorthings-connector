package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	cfgFlag := flag.String("config", "../configs/sampleconfig.json", "path of the config file")
	flag.Parse()
	cfg := *cfgFlag

	c, err := GetConfig(cfg)
	if err != nil {
		log.Fatal("config read error: ", err)
		return
	}

	start(c)

	go forever()
	select {} // block forever
}

func start(c Config) {
	publishChannel := make(chan *PublishMessage)
	publishBroker := CreatePubClient(c.PubBroker.Info.Host, c.Qos, c.ClientID, publishChannel, c.PubBroker.Info.Username, c.PubBroker.Info.Password)
	publishBroker.Start()

	subscribeBrokers := []MqttSubClient{}
	for _, sb := range c.SubBrokers {
		subCient := CreateSubClient(sb.Info.Host, c.Qos, sb.Streams, c.ClientID, publishChannel, sb.Info.Username, sb.Info.Password)
		subCient.Start()

		subscribeBrokers = append(subscribeBrokers, subCient)
	}
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}
