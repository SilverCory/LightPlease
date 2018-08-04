package LightPlease

import (
	"github.com/stianeikeland/go-rpio"
)

const PWMFreq = 510000

type IOOut struct {
	pinRed, pinGreen, pinBlue, pinWhite rpio.Pin
}

// Suggested pins
// 12, 13, 18, 19, 40, 41, 45
func NewIOOut(pinRed, pinGreen, pinBlue, pinWhite int) *IOOut {

	err := rpio.Open()
	if err != nil {
		return nil
	}

	ret := &IOOut{
		pinRed:   rpio.Pin(pinRed),
		pinGreen: rpio.Pin(pinGreen),
		pinBlue:  rpio.Pin(pinBlue),
		pinWhite: rpio.Pin(pinWhite),
	}

	ConfigurePWMPin(ret.pinRed)
	ConfigurePWMPin(ret.pinBlue)
	ConfigurePWMPin(ret.pinGreen)
	ConfigurePWMPin(ret.pinWhite)

	return ret
}

func ConfigurePWMPin(pin rpio.Pin) {
	rpio.PinMode(pin, rpio.Pwm)
	rpio.SetFreq(pin, PWMFreq)
	rpio.SetDutyCycle(pin, 0, 255)
}

func (i *IOOut) DisplayRGBW(R, G, B, W int16) {
	rpio.SetDutyCycle(i.pinRed, uint32(R), 255)
	rpio.SetDutyCycle(i.pinGreen, uint32(G), 255)
	rpio.SetDutyCycle(i.pinBlue, uint32(B), 255)
	rpio.SetDutyCycle(i.pinWhite, uint32(W), 255)
}

func (i *IOOut) Close() error {
	return rpio.Close()
}
