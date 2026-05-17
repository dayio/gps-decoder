package gps

import "fmt"

type NavDataState struct {
	BufferCount   int     // Counts up to 20ms
	Accumulator   float64 // Integrates the In-Phase (I) over 20ms
	ShiftRegister uint8   // 8 bit shift register (to find the preamble)
}

func ProcessNavigationData(prn int, I float64, state *NavDataState) {
	if prn != 15 {
		return
	}

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

		// Check for Preamble
		if state.ShiftRegister == 0x8B { // 0001011 (0x8B
			fmt.Printf("\n[PRN %02d] Preamble found (Normal) \n", prn)
		} else if state.ShiftRegister == 0x74 { // 01110100 (0x74)
			fmt.Printf("\n[PRN %02d] Preamble found (Inverted) \n", prn)
		}

		fmt.Printf("%d", bit)

		state.BufferCount = 0
		state.Accumulator = 0.0
	}
}
