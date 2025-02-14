package serial

import "math"

func roundTo3(x float64) float64 {
	return math.Floor(x*1000) / 1000
}
