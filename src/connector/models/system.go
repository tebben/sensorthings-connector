package models

type System interface {
	AddModule(ConnectorModule)
	GetModules() ([]ConnectorModule, error)
	GetConnectors() ([]Connector, error)
	GetConnector(id string) (Connector, error)
	GetEndpoints() []ConnectorEndpoint

	CreateConnector(connector *ConnectorBase) (Connector, error)
	PatchConnector(id string, connector *ConnectorBase) (Connector, error)
	DeleteConnector(id string) error

	SetConnectorState(id string, running bool) error

	Start()
}
