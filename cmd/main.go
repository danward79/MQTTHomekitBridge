package main

import (
	"fmt"

	bridge "github.com/danward79/MQTTHomekitBridge"
	"github.com/danward79/MQTTHomekitBridge/lightSensor"
	"github.com/danward79/MQTTHomekitBridge/logging"
	"github.com/danward79/MQTTHomekitBridge/temperatureSensor"
)

var (
	l *logging.Logger
)

func init() {
	l = logging.New("MQTTBridgeDeamon")
	l.Enable()

	if err := loadConfig(); err != nil {
		l.Fatal(fmt.Sprintf("Error reading config: %v", err))
	}
}

func main() {
	l.Message("Started")

	bridge := bridge.NewBridge(brokerIP, pinCode, "MQTTBridge", "MQTTBridge", "me!")

	bridge.AddServices(readConfig())

	bridge.Start()
}

// TODO: Add ability to read TOML config file.....
func readConfig() []bridge.BridgeableDevice {
	devs := []bridge.BridgeableDevice{}

	devs = append(devs, temperatureSensor.New("home/bedroom/temp", "Bedroom Temperature"))
	devs = append(devs, temperatureSensor.New("home/lounge/temp", "Lounge Temperature"))
	devs = append(devs, temperatureSensor.New("home/balcony/temp", "Balcony Temperature"))

	devs = append(devs, lightSensor.New("home/bedroom/light", "Bedroom Light"))
	devs = append(devs, lightSensor.New("home/lounge/light", "Lounge Light"))
	devs = append(devs, lightSensor.New("home/balcony/light", "Balcony Light"))

	return devs
}
