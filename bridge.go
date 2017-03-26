package bridge

import (
	"log"
	"os"
	"time"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	l "github.com/brutella/hc/log"
	"github.com/brutella/hc/service"
	"github.com/danward79/MQTTHomekitBridge/logging"
	"github.com/danward79/Thingamabob/clientService"
)

// AccessoryTypes stores a list of available accessory device types for the bridge
var AccessoryTypes = []string{"temperaturesensor", "lightsensor"}

// BridgeableDevice is a device that can be bridged between homekit and mqtt
type BridgeableDevice interface {
	Update([]byte) error
	Service() *service.Service
	Acc() *accessory.Accessory
	Topic() string
}

// Bridge stores details of a bridge
type Bridge struct {
	*accessory.Accessory

	deviceList map[string]BridgeableDevice

	mqttClient *clientService.Client
	logger     *logging.Logger
}

// New provides a new bridge device for accepting services.
func New(brokerIP, pin, name, model, manufacturer string) *Bridge {
	b := Bridge{
		Accessory: accessory.New(accessory.Info{
			Name:         name,
			SerialNumber: "000",
			Model:        model,
			Manufacturer: manufacturer,
		}, accessory.TypeBridge),

		deviceList: make(map[string]BridgeableDevice),

		mqttClient: clientService.New(brokerIP, "MQTTHomekitBridge", false),

		logger: logging.New("MQTTHomekitBridge"),
	}

	b.logger.Enable()
	b.logger.Message("New Bridge device created")

	return &b
}

// Start transport
func (b *Bridge) Start(devices *[]BridgeableDevice) {
	b.logger.Message("Starting")

	b.populateDeviceList(devices)

	b.subscribeTopics()

	accessories := createAccessoryList(b.deviceList)

	t, err := hc.NewIPTransport(hc.Config{}, b.Accessory, accessories...)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {

		t.Stop()

		time.Sleep(500 * time.Millisecond)

		os.Exit(1)
	})

	t.Start()
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

func createAccessoryList(devicelist map[string]BridgeableDevice) []*accessory.Accessory {
	var a []*accessory.Accessory

	for _, v := range devicelist {
		a = append(a, v.Acc())
	}

	return a
}

func (b *Bridge) populateDeviceList(devices *[]BridgeableDevice) {
	for _, d := range *devices {
		if _, ok := b.deviceList[d.Topic()]; !ok {
			b.logger.Message("Adding Accessory: TOPIC", d.Topic())
			b.deviceList[d.Topic()] = d
		}
	}
}

// EnableLogs detailed debug logs
func (b *Bridge) EnableLogs() {
	l.Debug.Enable()
}

// DisableLogs detailed debug logs
func (b *Bridge) DisableLogs() {
	l.Debug.Disable()
}
