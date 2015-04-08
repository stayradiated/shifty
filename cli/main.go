package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"github.com/stayradiated/shifty"
)

func main() {

	pin17, _ := OpenPinForOutput(rpi.GPIO17)
	pin27, _ := OpenPinForOutput(rpi.GPIO27)
	pin22, _ := OpenPinForOutput(rpi.GPIO22)

	s := &shifty.ShiftRegister{
		DataPin:  pin17,
		LatchPin: pin27,
		ClockPin: pin22,
	}

	for {
		leds := byte(0)
		s.Update(leds)
		time.Sleep(500 * time.Millisecond)

		for i := uint(0); i < 8; i++ {
			leds |= 1 << i
			s.Update(leds)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// OpenPinForOutput opens up a GPIO pin for output
func OpenPinForOutput(pinId int) (gpio.Pin, error) {

	// open pin
	pin, err := rpi.OpenPin(pinId, gpio.ModeOutput)
	if err != nil {
		return nil, err
	}

	// turn the pin off when we exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pin.Clear()
			pin.Close()
			os.Exit(0)
		}
	}()

	return pin, nil
}
