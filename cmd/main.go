package main

import (
	"context"
	"encoding/json"
	"github.com/borankux/sensors"
	"log"
	"time"
)

var handle sensors.MessageHandler = func(data []byte, sensorId int, deviceType int, parsed *sensors.ParserComplex) {
	log.Printf("%+v\n", parsed.Acceleration)
}

func main() {
	manager := sensors.NewManager(handle, context.Background())
	manager.SetStatusHandler(func(ms *sensors.ManagerStatus) {
		j, _ := json.Marshal(ms)
		log.Println(string(j))
	})
	manager.AddSensor(sensors.NewPassive(1, 1, "0.0.0.0", "00", 400, sensors.DeviceAcceleration, time.Second*5))
	manager.StartAllSensors()
	for {

	}
}
