package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	connectorErrors "github.com/tebben/sensorthings-connector/src/connector/errors"
	"github.com/tebben/sensorthings-connector/src/connector/models"
	"io/ioutil"
)

// HandleGetModules retrieves all modules from SensorThings Connector
func HandleGetModules(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	connector := *s
	handle := func() (interface{}, error) { return connector.GetModules() }
	HandleGetRequest(w, r, &handle)
}

// HandleGetConnectors retrieves all configured connectors
func HandleGetConnectors(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	connector := *s
	handle := func() (interface{}, error) { return connector.GetConnectors() }
	HandleGetRequest(w, r, &handle)
}

// HandlePostConnector handles a new created connector
func HandlePostConnector(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	system := *s
	byteData, _ := ioutil.ReadAll(r.Body)
	connector := &models.ConnectorBase{}
	err := json.Unmarshal(byteData, connector)
	if err != nil {
		sendError(w, connectorErrors.NewBadRequestError(errors.New("Unable to parse connector")))
	} else {
		if con, err := system.CreateConnector(connector); err != nil {
			sendError(w, connectorErrors.NewBadRequestError(err))
		} else {
			sendJSONResponse(w, http.StatusCreated, con)
		}
	}
}

// HandleGetConnectorById retrieves a connector by id
func HandleGetConnectorById(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	system := *s
	handle := func() (interface{}, error) { return system.GetConnector(ps.ByName("id")) }
	HandleGetRequest(w, r, &handle)
}

// HandleStartConnector start a connector by id
func HandleStartConnector(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	system := *s
	if err := system.SetConnectorState(ps.ByName("id"), true); err != nil {
		sendError(w, err)
	} else {
		sendJSONResponse(w, http.StatusOK, nil)
	}
}

// HandleStopConnector stops a connector by id
func HandleStopConnector(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	system := *s

	if err := system.SetConnectorState(ps.ByName("id"), false); err != nil {
		sendError(w, err)
	} else {
		sendJSONResponse(w, http.StatusOK, nil)
	}
}

// HandleDeleteConnector deletes a connector by id
func HandleDeleteConnector(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	system := *s
	if err := system.DeleteConnector(ps.ByName("id")); err != nil {
		sendError(w, err)
	} else {
		sendJSONResponse(w, http.StatusOK, nil)
	}
}

// HandlePatchConnector patches a connector by given id
func HandlePatchConnector(w http.ResponseWriter, r *http.Request, ps httprouter.Params, s *models.System) {
	system := *s
	byteData, _ := ioutil.ReadAll(r.Body)
	connector := &models.ConnectorBase{}
	err := json.Unmarshal(byteData, connector)
	if err != nil {
		sendError(w, connectorErrors.NewBadRequestError(errors.New("Unable to parse connector")))
	} else {
		if con, err := system.PatchConnector(ps.ByName("id"), connector); err != nil {
			sendError(w, err)
		} else {
			sendJSONResponse(w, http.StatusOK, con)
		}
	}
}

// handleGetRequest is the default function to handle incoming GET requests
func HandleGetRequest(w http.ResponseWriter, r *http.Request, h *func() (interface{}, error)) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	handler := *h
	if data, err := handler(); err != nil {
		sendError(w, err)
	} else {
		sendJSONResponse(w, http.StatusOK, data)
	}
}

// sendJSONResponse sends the desired message to the user
// the message will be marshalled into an indented JSON format
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	var b []byte
	var err error

	if data != nil {
		b, err = json.MarshalIndent(data, "", "   ")
		if err != nil {
			log.Printf("%v", err.Error())
		}
	}

	w.Write(b)
}

// sendError creates an ErrorResponse message and sets it to the user
// using SendJSONResponse
func sendError(w http.ResponseWriter, error error) {
	// Set te status code, default 500 for error, check if there is an ApiError an get
	// the status code
	var statusCode = http.StatusInternalServerError
	if error != nil {
		switch e := error.(type) {
		case connectorErrors.APIError:
			statusCode = e.GetHTTPErrorStatusCode()
			break
		}
	}

	statusText := http.StatusText(statusCode)
	errorResponse := models.ErrorResponse{
		Error: models.ErrorContent{
			StatusText: statusText,
			StatusCode: statusCode,
			Message:    error.Error(),
		},
	}

	sendJSONResponse(w, statusCode, errorResponse)
}
