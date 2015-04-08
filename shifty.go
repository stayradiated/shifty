package shifty

import (
	"sync"

	"github.com/davecheney/gpio"
)

type ShiftRegister struct {
	sync.Mutex

	LatchPin gpio.Pin
	DataPin  gpio.Pin
	ClockPin gpio.Pin

	state byte
	pins  []*Pin
}

func (s *ShiftRegister) Pin(index int) *Pin {

	if s.pins == nil {
		s.pins = make([]*Pin, 8)
	}

	if s.pins[index] != nil {
		return s.pins[index]
	}

	pin := &Pin{
		Index:    index,
		register: s,
	}

	s.pins[index] = pin

	return pin
}

func (s *ShiftRegister) AllPins() []*Pin {
	// make sure all pins are created
	for i := 0; i < 8; i++ {
		s.Pin(i)
	}

	return s.pins
}

func (s *ShiftRegister) SetBit(index int, state bool) {
	s.Lock()
	defer s.Unlock()

	if state {
		s.state |= 1 << byte(index)
	} else {
		s.state &= 1 ^ (1 << byte(index))
	}
	s.shiftOut()
}

func (s *ShiftRegister) GetBit(index int) bool {
	s.Lock()
	defer s.Unlock()

	return (s.state>>byte(index))&1 == 1
}

func (s *ShiftRegister) shiftOut() {
	s.LatchPin.Clear()

	for i := byte(0); i < 8; i++ {

		if (s.state>>i)&1 == 1 {
			s.DataPin.Set()
		} else {
			s.DataPin.Clear()
		}

		s.ClockPin.Set()
		// time.Sleep(20 * time.Millisecond)
		s.ClockPin.Clear()
		// time.Sleep(20 * time.Millisecond)

	}

	s.LatchPin.Set()
}

type Pin struct {
	Index    int
	register *ShiftRegister
}

func (p *Pin) Set() {
	p.register.SetBit(p.Index, true)
}

func (p *Pin) Clear() {
	p.register.SetBit(p.Index, false)
}

func (p *Pin) Get() bool {
	return p.register.GetBit(p.Index)
}
