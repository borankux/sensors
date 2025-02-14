package sensors

import (
	"fmt"
	"math"
	"time"
)

type Displacement struct {
	Displacement float64   `json:"displacement"`
	Timestamp    time.Time `json:"timestamp"`
	Address      string    `json:"address"`
}

func ParseDisplacement(msg []byte) (*Displacement, error) {
	if len(msg) < 7 {
		return nil, fmt.Errorf("invalid message format: expected at least 7 bytes, got %d", len(msg))
	}
	address := fmt.Sprintf("%02x", msg[0])
	raw := uint32(msg[3])<<24 | uint32(msg[4])<<16 | uint32(msg[5])<<8 | uint32(msg[6])
	d := math.Floor((float64(raw)/65536.0)*100) / 100

	return &Displacement{
		Displacement: d,
		Timestamp:    time.Now(),
		Address:      address,
	}, nil
}

func (d *Displacement) ToApiData() string {
	return fmt.Sprintf("d=%.2fmm", d.Displacement)
}
