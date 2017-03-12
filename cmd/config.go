package main

import (
	"flag"

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
	//logMessage(fmt.Sprintf("Broker: %s", brokerIP))

	return nil
}
