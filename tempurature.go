package sensors

import (
	"errors"
	"fmt"
	"time"
)

type TempRH struct {
	Temperature      float64   `json:"temperature"`
	RelativeHumidity float64   `json:"relative_humidity"`
	Timestamp        time.Time `json:"timestamp"`
	Address          string    `json:"address"`
}

func ParseTempRH(msg []byte) (*TempRH, error) {
	if len(msg) < 7 {
		return nil, errors.New("invalid message format: message too short")
	}

	address := fmt.Sprintf("%02x", msg[0])
	tempRaw := int16(uint16(msg[3])<<8 | uint16(msg[4]))
	temperature := float64(tempRaw) / 10.0
	rh := float64(uint16(msg[5])<<8|uint16(msg[6])) / 10.0

	return &TempRH{
		Temperature:      temperature,
		RelativeHumidity: rh,
		Timestamp:        time.Now(),
		Address:          address,
	}, nil
}

func (tr *TempRH) ToApiData() string {
	return fmt.Sprintf("temp=%.1f℃,rh=%.1f%%", tr.Temperature, tr.RelativeHumidity)
}
