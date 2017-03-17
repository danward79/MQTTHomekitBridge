package main

import (
	"flag"
	"fmt"

	bridge "github.com/danward79/MQTTHomekitBridge"
	"github.com/danward79/MQTTHomekitBridge/lightSensor"
	"github.com/danward79/MQTTHomekitBridge/temperatureSensor"
	"github.com/spf13/viper"
)

const (
	// PIN default pin code
	PIN = "00102003"
)

var (
	brokerIP, pinCode string
)

func loadConfig() error {

	broker := flag.String("b", "", `MQTT Broker address and port, 127.0.0.1:1883, excluding "tcp://"`)
	configPath := flag.String("c", "", "Config file location, default is working directory")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if *configPath != "" {
		viper.AddConfigPath(*configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetDefault("broker", "127.0.0.1:1883")
	if *broker != "" {
		viper.Set("broker", *broker)
	}
	brokerIP = viper.GetString("broker")

	viper.SetDefault("pin", PIN)
	pinCode = viper.GetString("pin")

	l.Message("Pin:", pinCode)
	l.Message("Broker:", brokerIP)

	return nil
}

func readConfigFile() []bridge.BridgeableDevice {
	devs := []bridge.BridgeableDevice{}

	for _, t := range bridge.AccessoryTypes {

		deviceKey := fmt.Sprintf("devices.%s", t)

		if viper.IsSet(deviceKey) {

			switch t {
			case "temperaturesensor":

				for _, v := range viper.Get(deviceKey).([]interface{}) {
					i := v.(map[string]interface{})
					devs = append(devs, temperatureSensor.New(i["topic"].(string), i["displayname"].(string)))
				}

			case "lightsensor":

				for _, v := range viper.Get(deviceKey).([]interface{}) {

					i := v.(map[string]interface{})
					devs = append(devs, lightSensor.New(i["topic"].(string), i["displayname"].(string)))
				}
			}

		}
	}

	return devs
}

// readConfig previously used to load hard-coded configuration.
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
