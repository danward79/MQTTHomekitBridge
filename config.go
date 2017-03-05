package main

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

var brokerIP string

func loadConfig() error {

	broker := flag.String("b", "", `MQTT Broker address and port, 127.0.0.1:1883, excluding "tcp://"`)
	configPath := flag.String("c", "", "Config file location, default is working directory")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if *configPath != "" {
		viper.AddConfigPath(*configPath)
	}

	viper.SetDefault("broker", "127.0.0.1:1883")
	if *broker != "" {
		viper.Set("broker", *broker)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	brokerIP = viper.GetString("broker")
	logMessage(fmt.Sprintf("Broker: %s", brokerIP))

	return nil
}

func readDeviceDetails(a *Accessories) {

	for _, t := range accessoryTypes {
		//log.Println("[MQTTHomekitBrige] Reading", t)

		if viper.IsSet(t) {
			for _, v := range viper.Get(t).([]interface{}) {

				v = v.(map[string]interface{})

				switch value := v.(type) {

				case interface{}:
					i := value.(map[string]interface{})

					logMessage(fmt.Sprintf("Found: %s, %v", t, i))

					accessoryDevice := map[string]string{
						"displayName":  i["displayname"].(string),
						"model":        i["model"].(string),
						"manufacturer": i["manufacturer"].(string),
						"topic":        i["topic"].(string),
						"min":          i["min"].(string),
						"max":          i["max"].(string),
						"inc":          i["inc"].(string),
						"serviceType":  t,
					}
					a.Add(accessoryDevice)

				default:
					logMessage(fmt.Sprintf("Unrecognised type during device read: %s, %s", reflect.TypeOf(value), value))
				}
			}
		}
	}
}
