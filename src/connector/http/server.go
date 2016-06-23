package http

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tebben/sensorthings-connector/src/connector/models"
)

// ConnectorHTTPServer is the type that contains all of the relevant information to set
// up the Connector HTTP Server
type ConnectorHTTPServer struct {
	system    *models.System
	host      string                     // Hostname for example "localhost:8081" or "192.168.1.14:8081"
	endpoints []models.ConnectorEndpoint // Configured endpoints for Connector HTTP
}

// CreateServer initialises a new Connector HTTPServer based on the given parameters
func CreateServer(system *models.System, host string, endpoints []models.ConnectorEndpoint) models.HTTPServer {
	return &ConnectorHTTPServer{
		system:    system,
		host:      host,
		endpoints: endpoints,
	}
}

// Start command to start the Connector HTTPServer
func (c *ConnectorHTTPServer) Start() {
	log.Printf("Started SensorThings Connector HTTP Server on %v", c.host)
	router := createRouter(c)
	httpError := http.ListenAndServe(c.host, router)

	if httpError != nil {
		log.Fatal(httpError)
		return
	}
}

// Stop command to stop the Connector HTTP server, currently not supported
func (c *ConnectorHTTPServer) Stop() {

}

func createRouter(c *ConnectorHTTPServer) *httprouter.Router {
	router := httprouter.New()
	for _, endpoint := range c.endpoints {
		ep := endpoint
		for _, op := range ep.GetOperations() {
			operation := op
			if operation.Handler == nil {
				continue
			}

			switch operation.OperationType {
			case models.HTTPOperationGet:
				{
					router.GET(operation.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						operation.Handler(w, r, p, c.system)
					})
				}
			case models.HTTPOperationPost:
				{
					router.POST(operation.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						operation.Handler(w, r, p, c.system)
					})
				}
			case models.HTTPOperationPatch:
				{
					router.PATCH(operation.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						operation.Handler(w, r, p, c.system)
					})
				}
			case models.HTTPOperationDelete:
				{
					router.DELETE(operation.Path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
						operation.Handler(w, r, p, c.system)
					})
				}
			}
		}
	}

	//router.ServeFiles("/app/*filepath", http.Dir("client/app"))
	return router
}
