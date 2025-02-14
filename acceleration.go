package sensors

import (
	"errors"
	"fmt"
	"time"
)

type Acceleration struct {
	AngularAccelX float64   `json:"angular_accel_x,omitempty"`
	AngularAccelY float64   `json:"angular_accel_y,omitempty"`
	AccelX        float64   `json:"accel_x,omitempty"`
	AccelY        float64   `json:"accel_y,omitempty"`
	AccelZ        float64   `json:"accel_z,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Address       string    `json:"address,omitempty"`
}

func decodeBCD(b byte) int {
	return int(b>>4)*10 + int(b&0x0F)
}

func parseBCD(signByte, intByte, fracByte byte) float64 {
	sign := 1.0
	if signByte != 0 {
		sign = -1.0
	}
	integerPart := decodeBCD(intByte)
	fraction := float64(decodeBCD(fracByte)) / 100.0
	return sign * (float64(integerPart) + fraction)
}

func ParseAcceleration(msg []byte) (*Acceleration, error) {
	if len(msg) < 27 {
		return nil, errors.New("invalid message format: message too short")
	}
	address := fmt.Sprintf("%02x", msg[2])
	angularAccelX := parseBCD(msg[8], msg[9], msg[10])
	angularAccelY := parseBCD(msg[12], msg[13], msg[14])
	accelX := parseBCD(msg[16], msg[17], msg[18]) / 10.0
	accelY := parseBCD(msg[20], msg[21], msg[22]) / 10.0
	accelZ := parseBCD(msg[24], msg[25], msg[26]) / 10.0
	return &Acceleration{
		AngularAccelX: roundTo3(angularAccelX),
		AngularAccelY: roundTo3(angularAccelY),
		AccelX:        roundTo3(accelX),
		AccelY:        roundTo3(accelY),
		AccelZ:        roundTo3(accelZ),
		Timestamp:     time.Now(),
		Address:       address,
	}, nil
}

func (a *Acceleration) ToApiData() string {
	return fmt.Sprintf("xaa=%.2f°,yaa=%.2f°,xa=%.2fg,ya=%.2fg,za=%.2fg",
		a.AngularAccelX, a.AngularAccelY, a.AccelX, a.AccelY, a.AccelZ)
}
