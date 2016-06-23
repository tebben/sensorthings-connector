package system

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/tebben/sensorthings-connector/src/connector/config"
	"github.com/tebben/sensorthings-connector/src/connector/database"
	connectorErrors "github.com/tebben/sensorthings-connector/src/connector/errors"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"github.com/tebben/sensorthings-connector/src/connector/mqtt"
	"github.com/tebben/sensorthings-connector/src/connector/rest"
)

type SensorThingsConnector struct {
	typeRegistry  map[string]reflect.Type
	connectors    map[string]models.Connector
	modules       []models.ConnectorModule
	restEndpoints []models.ConnectorEndpoint
	pubChannel    chan *models.PublishMessage
	pubClient     mqtt.MqttPubClient
	db            database.Database
}

// CreateSystem initialises a new SensorThings System
func CreateSystem(config config.Config) models.System {
	pubChan := make(chan *models.PublishMessage)
	pubClient := mqtt.CreatePubClient(config.PubBroker.Host, config.PubClient.Qos, config.PubClient.ClientID, pubChan, config.PubBroker.Username, config.PubBroker.Password)

	return &SensorThingsConnector{
		typeRegistry: make(map[string]reflect.Type, 0),
		connectors:   make(map[string]models.Connector, 0),
		pubChannel:   pubChan,
		pubClient:    pubClient,
		db:           database.Database{},
	}
}

// Start SensorThings connector, Start setups the modules and registers all module actions
func (sc *SensorThingsConnector) Start() {
	sc.restEndpoints = rest.CreateEndPoints()
	sc.pubClient.Start()
	// Load connectors from database
	sc.db.Open()
	connectors, err := sc.db.GetConnectors()
	if err != nil {
		log.Printf("%v", err.Error())
	} else {
		// Setup connector
		for idx, _ := range connectors {
			con := connectors[idx]
			if err = sc.setupConnector(con); err != nil {
				log.Printf("%v", err.Error())
				continue
			}

			sc.connectors[con.GetID()] = con
			con.Module.SettingsChanged(con.GetSettings())
			if con.GetIsRunning() {
				con.Start()
			}

			log.Printf("Connector loaded: %v", con.GetName())
		}
	}
}

// AddModule add a new module to SensorThings Connector
func (sc *SensorThingsConnector) AddModule(module models.ConnectorModule) {
	module.Setup()
	sc.modules = append(sc.modules, module)
	sc.typeRegistry[module.GetName()] = reflect.TypeOf(module)
}

// GetModules retrieves all current models added to SensorThings Connector
func (sc *SensorThingsConnector) GetModules() ([]models.ConnectorModule, error) {
	return sc.modules, nil
}

// GetConnectors retrieves all current created connectors
func (sc *SensorThingsConnector) GetConnectors() ([]models.Connector, error) {
	v := make([]models.Connector, 0, len(sc.connectors))

	for _, value := range sc.connectors {
		v = append(v, value)
	}

	return v, nil
}

// GetConnector retrieves a connector by id
func (sc *SensorThingsConnector) GetConnector(id string) (models.Connector, error) {
	if exist, err := sc.checkConnectorExist(id); !exist {
		return nil, err
	}

	return sc.connectors[id], nil
}

// GetEndpoints retrieves all REST endpoints defined for SensorThings Connector including module endpoints
func (sc *SensorThingsConnector) GetEndpoints() []models.ConnectorEndpoint {
	eps := make([]models.ConnectorEndpoint, 0)
	eps = append(eps, sc.restEndpoints...)

	return eps
}

// CreateConnector create a new connector based on given information and adds it to the database
func (sc *SensorThingsConnector) CreateConnector(connector *models.ConnectorBase) (models.Connector, error) {
	connector.ID = RandomString(8)
	if err := sc.db.InsertConnector(connector); err != nil {
		return nil, connectorErrors.NewRequestInternalServerError(err)
	}

	if err := sc.setupConnector(connector); err != nil {
		return nil, connectorErrors.NewRequestInternalServerError(err)
	}

	sc.connectors[connector.ID] = connector
	connector.Module.SettingsChanged(connector.GetSettings())
	log.Printf("Connector created: %v", connector.GetName())
	return connector, nil
}

// SetConnectorState sets the running state for a given module, returns an error if
// module not found
func (sc *SensorThingsConnector) SetConnectorState(id string, running bool) error {
	if exist, err := sc.checkConnectorExist(id); !exist {
		return err
	}

	c := sc.connectors[id]
	if running {
		c.Start()
	} else {
		c.Stop()
	}

	return sc.db.SaveConnectorState(id, running)
}

// PatchConnector updates a given Connector, user is unable to change id
func (sc *SensorThingsConnector) PatchConnector(id string, connector *models.ConnectorBase) (models.Connector, error) {
	if exist, err := sc.checkConnectorExist(id); !exist {
		return nil, err
	}

	connector.ID = id
	connector.Running = sc.connectors[id].GetIsRunning()

	if err := sc.setupConnector(connector); err != nil {
		return connector, connectorErrors.NewRequestInternalServerError(err)
	}

	if err := connector.GetModule().SettingsChanged(connector.GetSettings()); err != nil {
		return nil, connectorErrors.NewBadRequestError(err)
	}

	if err := sc.db.InsertConnector(connector); err != nil {
		return connector, connectorErrors.NewRequestInternalServerError(err)
	}

	sc.connectors[id].Stop()
	sc.connectors[id] = connector

	if connector.Running {
		connector.Start()
	}

	return connector, nil
}

// DeleteConnector stops the given connector if running and deletes it from the database
func (sc *SensorThingsConnector) DeleteConnector(id string) error {
	if exist, err := sc.checkConnectorExist(id); !exist {
		return err
	}

	if sc.connectors[id].GetIsRunning() {
		sc.connectors[id].Stop()
	}

	delete(sc.connectors, id)
	sc.db.DeleteConnector(id)

	return nil
}

// checkConnectorExist checks if there is a connector for the given id if not
// a HTTP RequestNotFound is returned
func (sc *SensorThingsConnector) checkConnectorExist(id string) (bool, error) {
	if _, ok := sc.connectors[id]; !ok {
		return false, connectorErrors.NewRequestNotFound(errors.New(fmt.Sprintf("Connector %s not found", id)))
	}

	return true, nil
}

// setupConnector creates a working connector from ConnectorBase by searching for the used module
// and instantiating the module using reflection, if the given module from ConnectorBase is not
// present an error will return
func (sc *SensorThingsConnector) setupConnector(connector *models.ConnectorBase) error {
	if t, ok := sc.typeRegistry[connector.GetModuleName()]; !ok {
		return errors.New(fmt.Sprintf("Error initialising %v, module: %v not found", connector.GetName(), connector.ModuleName))
	} else {
		newObjPtr := reflect.New(t.Elem())
		mod := newObjPtr.Interface().(models.ConnectorModule)
		mod.Setup()
		mod.SetPublishChannel(sc.pubChannel)
		connector.Module = mod
	}

	return nil
}
