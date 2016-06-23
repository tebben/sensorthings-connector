package netatmo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/exzz/netatmo-api-go"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"log"
	"time"
)

// NetatmoModule adds support for publishing Netatmo weather module readings
// to a SensorThings MQTT server.
type NetatmoModule struct {
	models.ConnectorModuleBase
	settings      NetatmoSettings
	fetchInterval time.Duration
	client        *netatmo.Client
	ticker        *time.Ticker
}

// NetatmoSettings contains information on Netatmo login and sensor reading to datastream mappings
type NetatmoSettings struct {
	ClientID      string        `json:"clientId"`
	ClientSecret  string        `json:"clientSecret"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	FetchInterval time.Duration `json:"fetchIntervalSeconds"`
	Mappings      []Mapping     `json:"mappings"`
}

type Mapping struct {
	ModuleID     string `json:"moduleId"`
	DataType     string `json:"dataType"`
	PublishTopic string `json:"publishTopic"`
}

// Setup initialised the module by setting some default values
func (nm *NetatmoModule) Setup() {
	nm.Name = "Netatmo"
	nm.Description = "Publish Netatmo readings to a SensorThings server"
	nm.fetchInterval = 10
}

// Start receiving Netatmo readings and publish it to a SensorThings server
func (nm *NetatmoModule) Start() {
	if len(nm.settings.ClientID) == 0 || len(nm.settings.ClientSecret) == 0 || len(nm.settings.Username) == 0 || len(nm.settings.Password) == 0 {
		log.Println("Incomplete settings for Netatmo module")
		return
	}

	nm.run()
}

// Stop receiving Netatmo readings
func (nm *NetatmoModule) Stop() {
	nm.ticker.Stop()
}

// SettingsChanged will try to parse and set NetatmoSettings from a json.RawMessage
func (nm *NetatmoModule) SettingsChanged(settings json.RawMessage) error {
	s := NetatmoSettings{}
	if err := json.Unmarshal(settings, &s); err != nil {
		return errors.New("Unable to read Netatmo Module settings")
	}

	nm.settings = s
	return nil
}

func (nm *NetatmoModule) run() {
	var err error
	nm.client, err = netatmo.NewClient(netatmo.Config{
		ClientID:     nm.settings.ClientID,
		ClientSecret: nm.settings.ClientSecret,
		Username:     nm.settings.Username,
		Password:     nm.settings.Password,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	interval := nm.settings.FetchInterval
	if interval == 0 {
		interval = nm.fetchInterval
	}

	// Get some readings at start
	go nm.getReadings()

	nm.ticker = time.NewTicker(time.Second * interval)
	go func() {
		for range nm.ticker.C {
			nm.getReadings()
		}
	}()
}

func (nm *NetatmoModule) getReadings() {
	dc, err := nm.client.GetDeviceCollection()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, station := range dc.Stations() {
			go nm.handleReadings(station.Modules())
		}
	}
}

// ToDo: Lesser for loops -> create mappings?
func (nm *NetatmoModule) handleReadings(modules []*netatmo.Device) {
	for _, module := range modules {
		for _, mapping := range nm.settings.Mappings {
			if mapping.ModuleID == module.ID {
				ts, data := module.Data()
				for dataType, value := range data {
					if mapping.DataType == dataType {
						pm := &models.PublishMessage{}
						pm.Topic = mapping.PublishTopic
						pm.Observation = &models.Observation{}
						pm.Observation.Result = value
						pm.Observation.PhenomenonTime = time.Unix(int64(ts), 0).Format(time.RFC3339Nano)

						nm.PublishChannel <- pm
					}
				}
			}
		}
	}
}
