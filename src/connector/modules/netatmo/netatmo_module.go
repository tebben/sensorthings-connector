package netatmo

import (
	"encoding/json"
	"errors"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"log"
)

type NetatmoModule struct {
	models.ConnectorModuleBase
	settings *NetatmoSettings
}

type NetatmoSettings struct {
}

func (nm *NetatmoModule) Setup() {
	nm.Name = "Netatmo"
	nm.Description = "Publish Netatmo readings to a SensorThings server"
}

func (nm *NetatmoModule) Start() {
	if nm.settings == nil {
		return
	}

	log.Println("START NETATMO")
}

func (nm *NetatmoModule) Stop() {
	log.Println("STOP NETATMO")
}

func (nm *NetatmoModule) SettingsChanged(settings json.RawMessage) error {
	s := &NetatmoSettings{}
	if err := json.Unmarshal(settings, s); err != nil {
		return errors.New("Unable to read Netatmo Module settings")
	}

	nm.settings = s

	return nil
}
