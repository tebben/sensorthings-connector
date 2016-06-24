package beeclear

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"net/http"
	"strings"
	"time"
)

// BeeClearModule adds support for publishing BeeClear readings
// to a SensorThings MQTT server.
type BeeClearModule struct {
	models.ConnectorModuleBase
	settings      BeeClearSettings
	fetchInterval time.Duration
	ticker        *time.Ticker
}

// BeeClearSettings contains information on BeeClear login and reading to datastream mappings
type BeeClearSettings struct {
	BeeClearHost  string        `json:"bcHost"`
	FetchInterval time.Duration `json:"fetchIntervalSeconds"`
	Mappings      []Mapping     `json:"mappings"`
}

// Mapping describes which value needs to published to what topic
type Mapping struct {
	DataType     string `json:"dataType"`     // for instance "u" for current usage or "g" for current generating gas
	PublishTopic string `json:"publishTopic"` // SensorThings MQTT topic to publish to
}

// Setup initialised the module by setting some default values
func (bc *BeeClearModule) Setup() {
	bc.Name = "BeeClear"
	bc.Description = "Publish BeeClear readings to a SensorThings server"
	bc.fetchInterval = 600 //Default to 600 seconds
}

// Start receiving BeeClear readings and publish it to a SensorThings server
func (bc *BeeClearModule) Start() {
	//ToDo: Check settings
	bc.run()
}

// Stop receiving BeeClear readings
func (bc *BeeClearModule) Stop() {
	bc.ticker.Stop()
}

// SettingsChanged will try to parse and set BeeClearSettings from a json.RawMessage
func (bc *BeeClearModule) SettingsChanged(settings json.RawMessage) error {
	s := BeeClearSettings{}
	if err := json.Unmarshal(settings, &s); err != nil {
		return errors.New("Unable to read BeeClear Module settings")
	}

	// Remove trailing slash from host
	if strings.HasSuffix(s.BeeClearHost, "/") {
		s.BeeClearHost = s.BeeClearHost[:len(s.BeeClearHost)-1]
	}

	bc.settings = s
	return nil
}

func (bc *BeeClearModule) run() {
	interval := bc.settings.FetchInterval
	if interval == 0 {
		interval = bc.fetchInterval
	}

	bc.ticker = time.NewTicker(time.Second * interval)
	go func() {
		for range bc.ticker.C {
			// ToDo retrieve readings
			// ToDo create PublishMessage
			// ToDo send publish message to channel: bc.PublishChannel <- publishMessage

			// Sample
			url := fmt.Sprintf("%s/bc_usage?date=1445554800&duration=168&period=24", bc.settings.BeeClearHost)
			bcUsage := make(map[string]int64)

			if err := getJson(url, &bcUsage); err == nil {
				for _, mapping := range bc.settings.Mappings {
					//check if param exist
					if _, ok := bcUsage[mapping.DataType]; !ok {
						continue
					}

					pm := &models.PublishMessage{}
					pm.Topic = mapping.PublishTopic
					pm.Observation = &models.Observation{}
					pm.Observation.Result = bcUsage[mapping.DataType]
					pm.Observation.PhenomenonTime = time.Unix(bcUsage["d"], 0).Format(time.RFC3339Nano)

					bc.PublishChannel <- pm
				}
			}
		}
	}()
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
