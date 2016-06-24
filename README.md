[![PyPI](https://img.shields.io/pypi/status/Django.svg?maxAge=2592000)]() [![apm](https://img.shields.io/apm/l/vim-mode.svg?maxAge=2592000)]() [![GoDoc](https://godoc.org/github.com/tebben/sensorthings-connector?status.svg)](https://godoc.org/github.com/tebben/sensorthings-connector) [![Go Report Card](https://goreportcard.com/badge/github.com/tebben/sensorthings-connector)](https://goreportcard.com/report/github.com/tebben/sensorthings-connector)
# sensorthings-connector
Middleware for publishing sensor readings to a SensorThings MQTT broker. New modules can be added easily to the sensorthings-connector to add support for a certain data source.

## sensorthings-connector config
When starting the sensorthings-connector a path to the config file needs to be supplied in the startup params for example:
./sensorthings-connector -config ./configs/sample.json

example of a config with comments
```
{
  "httpHost": "0.0.0.0:8081", // host where the HTTP server for the REST interface should run on
  "publishClient": { // definition of the client that will publish observations to a sensorthings MQTT broker
    "clientId": "stconnector", // id of the client
    "qos": 1, // quality of service
    "keepAlive": 300, // will set the amount of time (in seconds) that the client
                      // should wait before sending a PING request to the broker. This will
                      // allow the client to know that a connection has not been lost with the
                      // server.
    "pingTimeout": 20 // will set the amount of time (in seconds) that the client
                      // will wait after sending a PING request to the broker, before deciding
                      // that the connection has been lost
  },
  "publishBroker": { // definition of the sensorthings MQTT broker
      "host": "tcp://host:1883", // location of the broker to publish to
      "username": "", // supply username if needed
      "password": "" // supply password if needed
  }
}
```

## controlling the sensorthings-connector using REST
<u>Under scripts you can find a Postman file with example requests.</u>

<b>Get all modules</b>
```
GET: http://localhost:8081/Modules
STATUS: 200 OK
```

<b>Get all connectors</b>
```
GET: http://localhost:8081/Connectors
STATUS: 200 OK
```

<b>Get connector by id</b>
```
GET: http://localhost:8081/Connectors/{connectorID}
STATUS: 200 OK
```

<b>Create new connector</b>
```
GET: http://localhost:8081/Connectors
Body: {
         "name": "{connector name}",
         "description": "{connector description}",
         "module": "{module to use}",
         "settings": {
            {connector specific settings}
         }
       }
STATUS: 201 Created
```

<b>Update connector</b>
```
PATCH: http://localhost:8081/Connectors/{connectorID}
Body: {
         "name": "{connector name}",
         "description": "{connector description}",
         "module": "{name of the module to use}",
         "settings": {
            {connector specific settings}
         }
       }
STATUS: 200 OK
```

<b>Delete connector</b>
```
DELETE: http://localhost:8081/Connectors/{connectorID}
STATUS: 200 OK
```

<b>Start connector</b>
```
POST: http://localhost:8081/Connectors/{connectorID}/Start
STATUS: 200 OK
```

<b>Stop connector</b>
```
POST: http://localhost:8081/Connectors/{connectorID}/Stop
STATUS: 200 OK
```

## MODULES
### MQTT
MQTT can be used to connect an existing MQTT stream of sensor readings (using structured data) to the SensorThings broker.

Settings example when creating a connector using the MQTT module
```
"subBrokers": [
    {
      "host": "tcp://brokerhost:1883",
      "username": "",
      "password": "",
      "streams": [
        {
          "topicIn": "Test/1",
          "topicOut": "GOST/Datastreams(11)/Observations",
          "mapping": {
            "value": {
              "name": "result",
              "toFloat": true
            },
            "datetime": {
              "name": "phenomenonTime"
            }
          }
        }
      ]
    }
]
```
### Netatmo
Netatmo can be used to connect a Netatmo Weather Station to the SensorThings broker.

Settings example when creating a connector using the Netatmo module
```
{
    "name": "Netatmo VZ connector",
    "description": "Geodan VZ Netatmo readings connector",
    "module": "Netatmo",
    "settings": {
        "fetchIntervalSeconds": 600
        "clientId": "",
        "clientSecret": "",
        "username": "",
        "password": "",
        "mappings": [
            {
                "moduleId": "70:ee:50:03:65:d4",
                "dataType": "Temperature",
                "publishTopic": "GOST/Datastreams(3)/Observations"
            },
            {
                "moduleId": "70:ee:50:03:65:d4",
                "dataType": "Humidity",
                "publishTopic": "GOST/Datastreams(4)/Observations"
            },
            {
                "moduleId": "70:ee:50:03:65:d4",
                "dataType": "Pressure",
                "publishTopic": "GOST/Datastreams(5)/Observations"
            },
            {
                "moduleId": "70:ee:50:03:65:d4",
                "dataType": "CO2",
                "publishTopic": "GOST/Datastreams(6)/Observations"
            },
            {
                "moduleId": "70:ee:50:03:65:d4",
                "dataType": "Noise",
                "publishTopic": "GOST/Datastreams(7)/Observations"
            },
            {
                "moduleId": "02:00:00:03:5d:52",
                "dataType": "Temperature",
                "publishTopic": "GOST/Datastreams(8)/Observations"
            },
            {
                "moduleId": "02:00:00:03:5d:52",
                "dataType": "Humidity",
                "publishTopic": "GOST/Datastreams(9)/Observations"
            }
        ]
    }
}
```

Possible values for dataType:
Temperature, Humidity, Noise, Pressure, CO2,

fetchIntervalSeconds:
Not mandatory, defaults to 600 seconds (Unable to get faster readings from Netatmo API)

### BeeClear
ToDo
