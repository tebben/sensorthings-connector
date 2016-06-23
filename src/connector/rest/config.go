package rest

import "github.com/tebben/sensorthings-connector/src/connector/models"

// CreateEndPoints creates the pre-defined endpoint config, the config contains all endpoint info
func CreateEndPoints() []models.ConnectorEndpoint {
	endpoints := []models.ConnectorEndpoint{
		&Endpoint{
			Name: "Modules",
			Operations: []models.EndpointOperation{
				{models.HTTPOperationGet, "/Modules", HandleGetModules},
			},
		},
		&Endpoint{
			Name: "Connectors",
			Operations: []models.EndpointOperation{
				{models.HTTPOperationGet, "/Connectors", HandleGetConnectors},
				{models.HTTPOperationPost, "/Connectors", HandlePostConnector},
				{models.HTTPOperationGet, "/Connectors/:id", HandleGetConnectorById},
				{models.HTTPOperationPost, "/Connectors/:id/Start", HandleStartConnector},
				{models.HTTPOperationPost, "/Connectors/:id/Stop", HandleStopConnector},
				{models.HTTPOperationDelete, "/Connectors/:id", HandleDeleteConnector},
				{models.HTTPOperationPatch, "/Connectors/:id", HandlePatchConnector},
			},
		},
	}

	return endpoints
}
