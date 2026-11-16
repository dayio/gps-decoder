package gps

import "math"

type PLLState struct {
	DopplerFreq float64 // The corrected frequency in Hz (NCO)
	PhaseError  float64 // The error measured by the discriminator in radians
	GainP       float64 // Proportional gain
	GainI       float64 // Integral gain (frequency correction in Hz/rad)
}

type DLLState struct {
	CodePhase float64 // The exact position in the code (NCO) in chips
	CodeError float64 // The measured error
	Gain      float64 // The loop filter
}

func UpdatePLL(I, Q float64, state *PLLState) {

	if I == 0 {
		I = 1e-6 // Safety against division by zero
	}

	// DISCRIMINATOR : Costas Loop (2-quadrant arctan)
	state.PhaseError = math.Atan(Q / I)

	// Second-order loop filter
	gainP := state.GainP
	if gainP == 0 {
		gainP = 0.25
	}

	gainI := state.GainI
	if gainI == 0 {
		gainI = 2.0 // Hz per rad error per ms
	}

	// Integral path: update Doppler frequency
	state.DopplerFreq += state.PhaseError * gainI
}

func UpdateDLL(earlyPower, latePower, dopplerFreq float64, state *DLLState) {

	// DISCRIMINATOR : Normalized Early-Minus-Late
	sum := earlyPower + latePower

	if sum > 0 {
		state.CodeError = (earlyPower - latePower) / sum
	} else {
		state.CodeError = 0.0 // Safety against division by zero
	}

	// Code Doppler aiding from carrier tracking (1540 = 1575.42 MHz / 1.023 MHz)
	codeDoppler := (dopplerFreq / 1540.0) * 0.001

	// LOOP FILTER & NCO UPDATE: subtract error (negative feedback) and add Doppler aiding
	state.CodePhase += codeDoppler - state.CodeError*state.Gain

	state.CodePhase = math.Mod(state.CodePhase, 1023.0)

	if state.CodePhase < 0 {
		state.CodePhase += 1023.0
	}
}
