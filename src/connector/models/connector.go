package models

import "encoding/json"

// ConnectorModule describes all functions which will be called by the system
type ConnectorModule interface {
	GetName() string
	GetDescription() string
	SetPublishChannel(chan *PublishMessage)
	SettingsChanged(json.RawMessage) error
	Setup()
	Start()
	Stop()
}

// ConnectorModuleBase is a basic implementation of the ConnectorModule
type ConnectorModuleBase struct {
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	PublishChannel chan *PublishMessage `json:"-"`
}

// GetName returns the name of the module
func (mm *ConnectorModuleBase) GetName() string {
	return mm.Name
}

// GetDescription returns the description of the module
func (mm *ConnectorModuleBase) GetDescription() string {
	return mm.Description
}

// SetPublishChannel will be called by the system and passes in a channel where the module
// can pass PublishMessages to which then will be published to a MQTT broker
func (mm *ConnectorModuleBase) SetPublishChannel(channel chan *PublishMessage) {
	mm.PublishChannel = channel
}

// Connector defines a connector that can be created by the user, a connector instantiates a ConnectorModule
// so a module can be used multiple times for instance when you want to connect multiple Netatmo accounts
type Connector interface {
	GetID() string
	GetName() string
	GetDescription() string
	GetModule() ConnectorModule
	GetSettings() json.RawMessage
	GetIsRunning() bool

	Start()
	Stop()
}

// ConnectorBase is the default implementation of a Connector
type ConnectorBase struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ModuleName  string          `json:"module"`
	Running     bool            `json:"running"`
	Settings    json.RawMessage `json:"settings"`
	Module      ConnectorModule `json:"-"`
}

// GetID returns the id of the connector
func (c *ConnectorBase) GetID() string {
	return c.ID
}

// GetName returns the name of the connector, for instance: Tim's Netatmo readings
func (c *ConnectorBase) GetName() string {
	return c.Name
}

// GetModuleName returns the module name of the connector
func (c *ConnectorBase) GetModuleName() string {
	return c.ModuleName
}

// GetDescription returns the description of the Connector
func (c *ConnectorBase) GetDescription() string {
	return c.Description
}

// GetSettings returns the settings of the Connector
func (c *ConnectorBase) GetSettings() json.RawMessage {
	return c.Settings
}

// GetIsRunning returns if the connector is running or not
func (c *ConnectorBase) GetIsRunning() bool {
	return c.Running
}

// GetModule returns the instantiated ConnectorModule for the Connector
func (c *ConnectorBase) GetModule() ConnectorModule {
	return c.Module
}

// Start wil start running the connector
func (c *ConnectorBase) Start() {
	go c.GetModule().Start()
	c.Running = true
}

// Stop will stop the connector
func (c *ConnectorBase) Stop() {
	go c.GetModule().Stop()
	c.Running = false
}
