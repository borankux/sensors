package main

import (
	"context"
	"encoding/json"
	"github.com/borankux/sensors"
	"log"
	"time"
)

var handle serial.MessageHandler = func(data []byte, sensorId int, deviceType int, parsed *serial.ParserComplex) {
	log.Printf("%+v\n", parsed.Acceleration)
}

func main() {
	manager := serial.NewManager(handle, context.Background())
	manager.SetStatusHandler(func(ms *serial.ManagerStatus) {
		j, _ := json.Marshal(ms)
		log.Println(string(j))
	})
	manager.AddSensor(serial.NewPassive(1, 1, "0.0.0.0", "00", 400, serial.DeviceAcceleration, time.Second*5))
	manager.StartAllSensors()
	for {

	}
}
