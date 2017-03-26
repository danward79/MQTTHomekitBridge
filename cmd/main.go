package main

import (
	"fmt"

	bridge "github.com/danward79/MQTTHomekitBridge"
	"github.com/danward79/MQTTHomekitBridge/logging"
)

var (
	l *logging.Logger
)

func init() {
	l = logging.New("MQTTBridgeDeamon")
	l.Enable()

	l.Message("Started")

	if err := loadConfig(); err != nil {
		l.Fatal(fmt.Sprintf("Error reading config: %v", err))
	}

}

func main() {

	bridge := bridge.New(brokerIP, pinCode, bridgeName, "MQTTBridge", "N/A")
	//bridge.EnableLogs()
	bridge.Start(readConfigFile())
}
