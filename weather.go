package sensors

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type WeatherStation struct {
	Dn        float64   `json:"dn,omitempty"`
	Dm        float64   `json:"dm,omitempty"`
	Dx        float64   `json:"dx,omitempty"`
	Sn        float64   `json:"sn,omitempty"`
	Sm        float64   `json:"sm,omitempty"`
	Sx        float64   `json:"sx,omitempty"`
	Ta        float64   `json:"ta,omitempty"`
	Ua        float64   `json:"ua,omitempty"`
	Pa        float64   `json:"pa,omitempty"`
	Rc        float64   `json:"rc,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Address   string    `json:"address,omitempty"`
}

func ParseWeatherStation(message []byte) (*WeatherStation, error) {
	msg := string(message)
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.ReplaceAll(msg, "\n", "")
	pairs := strings.Split(msg, ",")
	var Dn, Dm, Dx, Sn, Sm, Sx, Ta, Ua, Pa, Rc float64
	zeroRegex := regexp.MustCompile(`^0*\.?0*$`)
	cleanRegex := regexp.MustCompile(`[^\d\.-]`)
	kvGet := func(kv []string, key string) float64 {
		if len(kv) < 2 {
			return 0
		}
		if kv[0] == key {
			value := kv[1]
			if zeroRegex.MatchString(value) {
				return 0
			}
			cleanedValue := cleanRegex.ReplaceAllString(value, "")
			num, err := strconv.ParseFloat(cleanedValue, 64)
			if err != nil {
				return 0
			}
			return num
		}
		return 0
	}

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)

		if len(kv) != 2 {
			continue
		}

		if val := kvGet(kv, "Dn"); val != 0 {
			Dn = val
		}

		if val := kvGet(kv, "Dm"); val != 0 {
			Dm = val
		}

		if val := kvGet(kv, "Dx"); val != 0 {
			Dx = val
		}

		if val := kvGet(kv, "Sn"); val != 0 {
			Sn = val
		}

		if val := kvGet(kv, "Sm"); val != 0 {
			Sm = val
		}

		if val := kvGet(kv, "Sx"); val != 0 {
			Sx = val
		}

		if val := kvGet(kv, "Ta"); val != 0 {
			Ta = val
		}

		if val := kvGet(kv, "Ua"); val != 0 {
			Ua = val
		}

		if val := kvGet(kv, "Pa"); val != 0 {
			Pa = val
		}

		if val := kvGet(kv, "Rc"); val != 0 {
			Rc = val
		}
	}

	ws := &WeatherStation{
		Dn:        Dn,
		Dm:        Dm,
		Dx:        Dx,
		Sn:        Sn,
		Sm:        Sm,
		Sx:        Sx,
		Ta:        Ta,
		Ua:        Ua,
		Pa:        Pa,
		Rc:        Rc,
		Timestamp: time.Now(),
		Address:   "00",
	}

	return ws, nil
}

func (ws *WeatherStation) ToApiData() string {
	return fmt.Sprintf("Dn=%.1fD,Dm=%.1fD,Dx=%.1fD,Sn=%.1fM,Sm=%.1fM,Sx=%.1fM,Ta=%.1fC,Ua=%.1fP,Pa=%.1fH,Rc=%.2fM",
		ws.Dn, ws.Dm, ws.Dx, ws.Sn, ws.Sm, ws.Sx, ws.Ta, ws.Ua, ws.Pa, ws.Rc)
}
