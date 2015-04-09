package shifty

import (
	"fmt"
	"strconv"
	"sync"
)

type Pin interface {
	Set()
	Clear()
	Get() bool
}

type ShiftRegister struct {
	sync.Mutex

	LatchPin Pin
	DataPin  Pin
	ClockPin Pin
	MaxPins  uint

	state uint
	pins  []Pin
}

func (s *ShiftRegister) Pin(index uint) Pin {

	if s.pins == nil {
		s.pins = make([]Pin, s.MaxPins)
	}

	if s.pins[index] != nil {
		return s.pins[index]
	}

	pin := &ShiftRegisterPin{
		Index:    index,
		register: s,
	}

	s.pins[index] = pin

	return pin
}

func (s *ShiftRegister) AllPins() []Pin {
	// make sure all pins are created
	for i := uint(0); i < s.MaxPins; i++ {
		s.Pin(i)
	}

	return s.pins
}

func (s *ShiftRegister) SetBit(index uint, state bool) {
	s.Lock()
	defer s.Unlock()

	if state {
		s.state |= 1 << index
	} else {
		s.state &= ^(1 << index)
	}

	s.shiftOut()
}

func (s *ShiftRegister) GetBit(index uint) bool {
	s.Lock()
	defer s.Unlock()

	return (s.state>>index)&1 > 0
}

func (s *ShiftRegister) shiftOut() {
	fmt.Println(strconv.FormatInt(int64(s.state), 2))

	s.LatchPin.Clear()

	for i := uint(0); i < s.MaxPins; i++ {

		if (s.state>>i)&1 > 0 {
			s.DataPin.Set()
		} else {
			s.DataPin.Clear()
		}

		s.ClockPin.Set()
		s.ClockPin.Clear()
	}

	s.LatchPin.Set()
}

type ShiftRegisterPin struct {
	Index    uint
	register *ShiftRegister
}

func (p *ShiftRegisterPin) Set() {
	p.register.SetBit(p.Index, true)
}

func (p *ShiftRegisterPin) Clear() {
	p.register.SetBit(p.Index, false)
}

func (p *ShiftRegisterPin) Get() bool {
	return p.register.GetBit(p.Index)
}
