package gps

import (
	"time"

	"github.com/dayio/gps-decoder/internal/events"
)

type NavDataState struct {
	BufferCount   int     // Counts up to 20ms
	Accumulator   float64 // Integrates the In-Phase (I) over 20ms
	ShiftRegister uint8   // 8 bit shift register (to find the preamble)
}

func ProcessNavigationData(prn int, I float64, state *NavDataState, bus *events.Bus, store *NavStore) {

	state.Accumulator += I
	state.BufferCount++

	if state.BufferCount >= 20 {

		// Determine the bit
		bit := uint8(0)
		if state.Accumulator < 0 {
			bit = 1
		}

		// Push the bit into shift register
		state.ShiftRegister = (state.ShiftRegister << 1) | bit

		if state.ShiftRegister == 0x8B {
			bus.Publish(events.Event{
				Timestamp: time.Now(),
				PRN:       prn,
				Type:      events.EventPreambleFound,
				Message:   "Normal preamble detected",
				Data:      events.PreambleNormal,
			})
		} else if state.ShiftRegister == 0x74 {
			bus.Publish(events.Event{
				Timestamp: time.Now(),
				PRN:       prn,
				Type:      events.EventPreambleFound,
				Message:   "Inverted preamble detected (180° phase flip)",
				Data:      events.PreambleInverted,
			})
		}

		// fmt.Printf("%d", bit)

		state.BufferCount = 0
		state.Accumulator = 0.0
	}
}
