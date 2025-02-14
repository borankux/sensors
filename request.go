package sensors

import (
	"encoding/hex"
	"errors"
)

func GetRequest(deviceType int, address string) ([]byte, error) {
	if deviceType == DeviceWeather {
		return []byte("0R0\r\n"), nil
	}
	if deviceType == DeviceTempRH || deviceType == DeviceDisplacement {
		addr, err := hex.DecodeString(address)
		if err != nil {
			return nil, err
		}
		if len(addr) == 0 {
			return nil, errors.New("invalid address")
		}
		return AppendCRC16([]byte{addr[0], 0x03, 0x00, 0x00, 0x00, 0x02}), nil
	}
	return nil, errors.New("invalid device type")
}
