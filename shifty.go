package shifty

import (
	"time"

	"github.com/davecheney/gpio"
)

type ShiftRegister struct {
	LatchPin gpio.Pin
	DataPin  gpio.Pin
	ClockPin gpio.Pin
}

func (s *ShiftRegister) Update(data byte) {
	s.LatchPin.Clear()

	s.ShiftOut(data)

	s.LatchPin.Set()
}

func (s *ShiftRegister) ShiftOut(data byte) {

	for i := byte(0); i < 8; i++ {

		if (data>>i)&1 == 1 {
			s.DataPin.Set()
		} else {
			s.DataPin.Clear()
		}

		s.ClockPin.Set()
		time.Sleep(20 * time.Millisecond)
		s.ClockPin.Clear()
		time.Sleep(20 * time.Millisecond)

	}

}
