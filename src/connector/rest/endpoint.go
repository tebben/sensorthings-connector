package rest

import "github.com/tebben/sensorthings-connector/src/connector/models"

// Endpoint contains all information for creating and handling an endpoint.
// Endpoint can be marshalled to JSON for returning endpoint information requested by the user
type Endpoint struct {
	Name       string                     `json:"name"` // Name of the endpoint
	Operations []models.EndpointOperation `json:"operations"`
}

// GetName returns the endpoint name
func (e *Endpoint) GetName() string {
	return e.Name
}

// GetOperations returns all operations for this endpoint such as GET, POST
func (e *Endpoint) GetOperations() []models.EndpointOperation {
	return e.Operations
}
