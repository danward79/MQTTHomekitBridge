package lightSensor

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

// CurrentAmbientLightLevel ...
type CurrentAmbientLightLevel struct {
	*characteristic.Float
}

// CustomLightSensor ...
type CustomLightSensor struct {
	*service.Service

	CurrentAmbientLightLevel *CurrentAmbientLightLevel
}

func newLightSensor() *CustomLightSensor {
	svc := CustomLightSensor{}
	svc.Service = service.New(service.TypeLightSensor)

	svc.CurrentAmbientLightLevel = newCurrentAmbientLightLevel()
	svc.AddCharacteristic(svc.CurrentAmbientLightLevel.Characteristic)

	return &svc
}

func newCurrentAmbientLightLevel() *CurrentAmbientLightLevel {
	char := characteristic.NewFloat(characteristic.TypeCurrentAmbientLightLevel)
	char.Format = characteristic.FormatFloat
	char.Perms = []string{characteristic.PermRead, characteristic.PermEvents}
	char.SetMinValue(0)
	char.SetMaxValue(100)
	char.SetStepValue(1)
	char.SetValue(0)
	char.Unit = characteristic.UnitPercentage //TODO: Fix units. Lux still showing...

	return &CurrentAmbientLightLevel{char}
}
