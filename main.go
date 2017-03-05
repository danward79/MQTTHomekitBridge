package main

import (
	"log"
)

func init() {
	logMessage("Started")

	if err := loadConfig(); err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}
}

func main() {
	accessories := New(brokerIP)
	readDeviceDetails(accessories)
	accessories.InitSensors()

	accessories.Run()

}

func logMessage(s string) {
	log.Printf("[MQTTHomekitBrige] %s\n", s)
}
