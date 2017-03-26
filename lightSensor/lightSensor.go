package lightSensor

import (
	"strconv"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

// LightSensor stores the details of a device that is being bridged
type LightSensor struct {
	*accessory.Accessory
	*CustomLightSensor

	topic string
}

// New returns a new LightSensor device
func New(topic, displayName string) *LightSensor {

	l := LightSensor{
		Accessory:         accessory.New(accessory.Info{Name: displayName}, accessory.TypeSensor), // TODO: Check TypeSensor
		CustomLightSensor: newLightSensor(),

		topic: topic,
	}

	l.CustomLightSensor.CurrentAmbientLightLevel.SetMinValue(0)
	l.CustomLightSensor.CurrentAmbientLightLevel.SetMaxValue(100)

	l.AddService(l.CustomLightSensor.Service)

	return &l
}

// Update a LightSensor device with a new received value.
func (l *LightSensor) Update(value []byte) error {
	v := parseFloat64(string(value), 0)
	l.CurrentAmbientLightLevel.SetValue(v)
	return nil
}

// Service returns the service
func (l *LightSensor) Service() *service.Service {
	return l.CustomLightSensor.Service
}

// Acc returns the accessory
func (l *LightSensor) Acc() *accessory.Accessory {
	return l.Accessory
}

// Topic returns the topic
func (l *LightSensor) Topic() string {
	return l.topic
}

func parseFloat64(s string, defaultValue float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return f
}
