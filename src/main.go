package main

import (
	"flag"
	"log"

	"github.com/tebben/sensorthings-connector/src/connector/config"
	"github.com/tebben/sensorthings-connector/src/connector/http"
	"github.com/tebben/sensorthings-connector/src/connector/modules/mqtt"
	"github.com/tebben/sensorthings-connector/src/connector/modules/netatmo"
	"github.com/tebben/sensorthings-connector/src/connector/system"
)

func main() {
	cfgFlag := flag.String("config", "configs/sampleconfig.json", "path of the config file")
	flag.Parse()
	cfg := *cfgFlag

	c, err := config.GetConfig(cfg)
	if err != nil {
		log.Fatal("config read error: ", err)
		return
	}

	start(c)
}

func start(c config.Config) {
	system := system.CreateSystem(c)

	//---ADD MODULES HERE---//
	system.AddModule(&mqtt.MQTTModule{})
	system.AddModule(&netatmo.NetatmoModule{})
	//----------------------//

	system.Start()

	connectorServer := http.CreateServer(&system, c.HttpHost, system.GetEndpoints())
	connectorServer.Start()
}
