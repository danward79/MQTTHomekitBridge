package bridge

import (
	"log"
	"os"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/danward79/MQTTHomekitBridge/logging"
	"github.com/danward79/Thingamabob/clientService"
)

// AccessoryDevices stores a list of available accessory device types for the bridge
var AccessoryTypes = []string{"temperaturesensor", "lightsensor"}

// BridgeableDevice is a device that can be bridged between homekit and mqtt
type BridgeableDevice interface {
	Update([]byte) error
	Service() *service.Service
	Topic() string
}

// Bridge stores details of a bridge
type Bridge struct {
	*accessory.Accessory

	transport       hc.Transport
	transportConfig hc.Config
	deviceList      map[string]BridgeableDevice

	mqttClient *clientService.Client
	logger     *logging.Logger
}

// NewBridge provides a new bridge device for accepting services.
func NewBridge(brokerIP, pin, name, model, manufacturer string) *Bridge {
	b := Bridge{
		Accessory: accessory.New(accessory.Info{
			Name:         name,
			Model:        model,
			Manufacturer: manufacturer,
		}, accessory.TypeBridge),
		transportConfig: hc.Config{Pin: pin},
		deviceList:      make(map[string]BridgeableDevice),
		mqttClient:      clientService.New(brokerIP, "MQTTHomekitBridge", false),
		logger:          logging.New("MQTTHomekitBridge"),
	}

	b.logger.Enable()
	b.logger.Message("New Bridge device created")
	return &b
}

// AddServices add BridgeableDevice to the Bridge
func (b *Bridge) AddServices(devices []BridgeableDevice) {

	for _, d := range devices {
		b.AddService(d.Service())

		if _, ok := b.deviceList[d.Topic()]; !ok {
			b.logger.Message("Adding service: TOPIC", d.Topic())
			b.deviceList[d.Topic()] = d
		}
	}
}

// Start transport
func (b *Bridge) Start() {
	b.logger.Message("Starting")

	b.subscribeTopics()

	var err error
	b.transport, err = hc.NewIPTransport(b.transportConfig, b.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		if b.transport != nil {
			b.transport.Stop()
		}

		os.Exit(1)
	})

	b.transport.Start()
}

func (b *Bridge) subscribeTopics() chan clientService.Message {

	var topics []string
	for k := range b.deviceList {
		topics = append(topics, k)
	}

	f := b.mqttClient.SubscribeMultiple(topics)

	b.logger.Message("Subscribing topics")

	go b.watch(f)
	return f
}

func (b *Bridge) watch(f <-chan clientService.Message) {
	b.logger.Message("Watching topics")

	for msg := range f {
		b.deviceList[msg.Topic()].Update(msg.Payload())
	}

}
