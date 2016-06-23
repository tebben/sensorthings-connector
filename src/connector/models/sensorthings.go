package models

// PublishMessage is used to publish a message trough the PubClient.
// When a subscription client receives a message it will be transformed into a PublishMessage
// and send trough the publish channel where the publish broker will pick up the message
type PublishMessage struct {
	Topic       string       `json:"topic"`
	Observation *Observation `json:"observation"`
}

// Observation in SensorThings represents a single Sensor reading of an ObservedProperty. A physical device, a Sensor, sends
// Observations to a specified Datastream. An Observation requires a FeatureOfInterest entity, if none is provided in the request,
// the Location of the Thing associated with the Datastream, will be assigned to the new Observation as the FeaturOfInterest.
type Observation struct {
	PhenomenonTime    string                 `json:"phenomenonTime,omitempty"`
	Result            interface{}            `json:"result,omitempty"`
	ResultTime        string                 `json:"resultTime,omitempty"`
	ResultQuality     string                 `json:"resultQuality,omitempty"`
	ValidTime         string                 `json:"validTime,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	FeatureOfInterest *FeatureOfInterest     `json:"FeatureOfInterest,omitempty"`
}

// FeatureOfInterest in SensorThings represents the phenomena an Observation is detecting. In some cases a FeatureOfInterest
// can be the Location of the Sensor and therefore of the Observation. A FeatureOfInterest is linked to a single Observation
type FeatureOfInterest struct {
	NavSelf      string                 `json:"@iot.selfLink,omitempty"`
	Description  string                 `json:"description,omitempty"`
	EncodingType string                 `json:"encodingtype,omitempty"`
	Feature      map[string]interface{} `json:"feature,omitempty"`
}
