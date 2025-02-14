package serial

import (
	"encoding/json"
	"time"
)

type ApiSchema struct {
	DeviceId        int       `json:"device_id"`
	Type            int       `json:"type"`
	DeviceTimestamp time.Time `json:"device_timestamp"`
	Data            string    `json:"data"`
	Project         string    `json:"project"`
	Name            string    `json:"name"`
}

func (api ApiSchema) ToJson() (string, error) {
	marshaled, err := json.Marshal(api)
	if err != nil {
		return "", err
	}
	return string(marshaled), nil
}

type Parser interface {
	ToApiData() string
}

func BuildApiRequestBody(projectName string, deviceId int, deviceType int, deviceName string, parser Parser) *ApiSchema {
	return &ApiSchema{
		DeviceId:        deviceId,
		Type:            deviceType,
		DeviceTimestamp: time.Now(),
		Data:            parser.ToApiData(),
		Project:         projectName,
		Name:            deviceName,
	}
}
