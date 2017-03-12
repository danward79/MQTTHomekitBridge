package bridge

import (
	"fmt"
	"log"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/danward79/Thingamabob/clientService"
)

// BridgeableDevice is a device that can be bridged between homekit and mqtt
type BridgeableDevice interface {
	Update([]byte) error
	Service() *service.Service
	Topic() string
}

// Bridge stores details of a bridge
type Bridge struct {
	*accessory.Accessory

	transport  hc.Transport
	deviceList map[string]BridgeableDevice

	mqttClient *clientService.Client
}

// NewBridge provides a new bridge device for accepting services.
func NewBridge(brokerIP string) *Bridge {
	b := Bridge{
		Accessory: accessory.New(accessory.Info{
			Name: "MQTTBridge",
			//Model:        "MQTTBridge",
			//Manufacturer: "Dan!",
		}, accessory.TypeBridge),
		deviceList: make(map[string]BridgeableDevice),
		mqttClient: clientService.New(brokerIP, "MQTTHomekitBridge", false),
	}

	return &b
}

// AddServices add BridgeableDevice to the Bridge
func (b *Bridge) AddServices(devices []BridgeableDevice) {

	for _, d := range devices {
		b.AddService(d.Service())

		if _, ok := b.deviceList[d.Topic()]; !ok {
			fmt.Println("TOPIC", d.Topic())
			b.deviceList[d.Topic()] = d
		}
	}
}

// Start transport
func (b *Bridge) Start() {
	fmt.Println("subscribe ")
	//msgFeed :=
	b.subscribeTopics()
	fmt.Println("subscribe done")

	var err error
	b.transport, err = hc.NewIPTransport(hc.Config{}, b.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Bla")
	//go watch(msgFeed)

	//go b.transport.Start()
	b.transport.Start()
}

func (b *Bridge) subscribeTopics() chan clientService.Message {

	var topics []string
	for k := range b.deviceList {
		topics = append(topics, k)
	}

	f := b.mqttClient.SubscribeMultiple(topics)
	//f := b.mqttClient.Subscribe("home/bedroom/temp")
	go b.watch(f)
	return f
}

func (b *Bridge) watch(f <-chan clientService.Message) {
	fmt.Println("WATCHING")
	for msg := range f {
		b.deviceList[msg.Topic()].Update(msg.Payload())
	}

}
