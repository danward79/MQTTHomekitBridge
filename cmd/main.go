package main

import (
	bridge "github.com/danward79/MQTTHomekitBridge"
	"github.com/danward79/MQTTHomekitBridge/lightSensor"
	"github.com/danward79/MQTTHomekitBridge/logging"
	"github.com/danward79/MQTTHomekitBridge/temperatureSensor"
)

func main() {
	l := logging.New("MQTTBridgeDeamon")
	l.Message("Started")

	bridge := bridge.NewBridge("localhost:1883")

	bridge.AddServices(readConfig())

	bridge.Start()

}

func readConfig() []bridge.BridgeableDevice {
	devs := []bridge.BridgeableDevice{}

	devs = append(devs, temperatureSensor.New("home/bedroom/temp", "Bedroom Temperature"))
	devs = append(devs, temperatureSensor.New("home/lounge/temp", "Lounge Temperature"))
	devs = append(devs, temperatureSensor.New("home/balcony/temp", "Balcony Temperature"))

	devs = append(devs, lightSensor.New("home/bedroom/light", "Bedroom Light"))
	//devs = append(devs, lightSensor.New("home/lounge/light", "Lounge Light"))
	//devs = append(devs, lightSensor.New("home/balcony/light", "Balcony Light"))

	return devs
}
