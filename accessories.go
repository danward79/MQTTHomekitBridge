package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/danward79/Thingamabob/clientService"
)

var accessoryTypes = []string{"temperaturesensor"} //, "lightsensor"

// Accessories ...
type Accessories struct {
	mqttClient *clientService.Client
	acc        map[string]accessoryParams
	transport  hc.Transport
	msgFeed    chan clientService.Message
}

type accessoryParams struct {
	info        accessory.Info
	topic       string
	min         float64
	max         float64
	inc         float64
	serviceType string
	acc         *accessory.Accessory
	initialised bool
}

// New return an empty list of accessiores
func New(b string) *Accessories {
	return &Accessories{
		mqttClient: clientService.New(brokerIP, "MQTTHomekitBridge", false),
		acc:        make(map[string]accessoryParams),
	}
}

// Add an accessory to the list of accessiores available
func (a *Accessories) Add(p map[string]string) {

	ap := accessoryParams{
		info: accessory.Info{
			Name:         p["displayname"],
			Model:        p["model"],
			Manufacturer: p["manufacturer"],
			SerialNumber: p["serial"],
		},
		topic:       p["topic"],
		min:         parseFloat64(p["min"], -40),
		max:         parseFloat64(p["max"], 100),
		inc:         parseFloat64(p["inc"], 0.1),
		serviceType: p["service"],
		initialised: false,
	}
	//a.acc = append(a.acc, ap)
	a.acc[p["topic"]] = ap
}

// InitSensors ...
func (a *Accessories) InitSensors() {

	var accessoryArray []*accessory.Accessory
	var topics []string

	for _, v := range a.acc {
		topics = append(topics, v.topic)
		v.acc = accessory.NewTemperatureSensor(v.info, 0, v.min, v.max, v.inc).Accessory
		accessoryArray = append(accessoryArray, v.acc)
	}
	a.msgFeed = a.mqttClient.SubscribeMultiple(topics)

	var err error
	a.transport, err = hc.NewIPTransport(hc.Config{}, accessoryArray[0], accessoryArray[1:]...)
	if err != nil {
		log.Fatal(err)
	}
}

// Run ...
func (a *Accessories) Run() {

	hc.OnTermination(func() {
		a.transport.Stop()
	})

	go func() {
		for {
			msg := <-a.msgFeed
			fmt.Println(msg)
			fmt.Println("map", a.acc[msg.Topic()])
			fmt.Println("acc info", a.acc[msg.Topic()].acc.Info)
			if a.acc[msg.Topic()].acc != nil {
				//a.TempSensor.CurrentTemperature.SetValue(v)
				fmt.Println("ACC", a.acc[msg.Topic()].acc)
			}
		}
	}()

	a.start()

}

func (a *Accessories) start() {
	a.transport.Start()
}

func parseFloat64(s string, deFault float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return deFault
	}
	return f
}
