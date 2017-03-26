package temperatureSensor

import (
	"strconv"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

// TemperatureSensor stores the details of a device that is being bridged
type TemperatureSensor struct {
	*accessory.Accessory
	*service.TemperatureSensor

	topic string
}

// New returns a new TemperatureSensor device
func New(topic string, displayName string) *TemperatureSensor {

	t := TemperatureSensor{
		Accessory:         accessory.New(accessory.Info{Name: displayName}, accessory.TypeThermostat),
		TemperatureSensor: service.NewTemperatureSensor(),

		topic: topic,
	}

	t.TemperatureSensor.CurrentTemperature.SetMinValue(-40)
	t.TemperatureSensor.CurrentTemperature.SetMaxValue(100)

	t.AddService(t.TemperatureSensor.Service)

	return &t
}

// Update a TemperatureSensor device with a new received value.
func (t *TemperatureSensor) Update(value []byte) error {
	v := parseFloat64(string(value), 0)
	t.CurrentTemperature.SetValue(v)
	return nil
}

// Service returns the service
func (t *TemperatureSensor) Service() *service.Service {
	return t.TemperatureSensor.Service
}

// Acc returns the accessory
func (t *TemperatureSensor) Acc() *accessory.Accessory {
	return t.Accessory
}

// Topic returns the topic
func (t *TemperatureSensor) Topic() string {
	return t.topic
}

func parseFloat64(s string, defaultValue float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}

	return f
}
