package serial

import "fmt"

type ParserComplex struct {
	Acceleration   *Acceleration
	Displacement   *Displacement
	TempRH         *TempRH
	WeatherStation *WeatherStation
}

func Parse(data []byte, dataType int) (*ParserComplex, error) {
	if dataType == DeviceAcceleration {
		parsed, err := ParseAcceleration(data)
		if err != nil {
			return nil, err
		}
		return &ParserComplex{
			Acceleration: parsed,
		}, nil
	}

	if dataType == DeviceDisplacement {
		parsed, err := ParseDisplacement(data)
		if err != nil {
			return nil, err
		}
		return &ParserComplex{
			Displacement: parsed,
		}, nil
	}

	if dataType == DeviceTempRH {
		parsed, err := ParseTempRH(data)
		if err != nil {
			return nil, err
		}
		return &ParserComplex{
			TempRH: parsed,
		}, nil
	}

	if dataType == DeviceWeather {
		parsed, err := ParseWeatherStation(data)
		if err != nil {
			return nil, err
		}
		return &ParserComplex{
			WeatherStation: parsed,
		}, nil
	}

	return nil, fmt.Errorf("unsupported data type %d", dataType)
}
