package gps

import "math"

type PLLState struct {
	DopplerFreq float64 // The corrected frequency (NCO)
	PhaseError  float64 // The error measured by the discriminator
	Gain        float64 // The loop filter
}

type DLLState struct {
	CodePhase float64 // The exact position in the code (NCO)
	CodeError float64 // The measured error
	Gain      float64 // The loop filter
}

func UpdatePLL(I, Q float64, state *PLLState) {

	// DISCRIMINATOR : Costas Loop
	state.PhaseError = math.Atan2(Q, I)

	// LOOP FILTER & NCO UPDATE
	state.DopplerFreq += state.PhaseError * state.Gain
}

func UpdateDLL(earlyPower, latePower float64, state *DLLState) {

	// DISCRIMINATOR : Normalized Early-Minus-Late
	sum := earlyPower + latePower

	if sum > 0 {
		state.CodeError = (earlyPower - latePower) / sum
	} else {
		state.CodeError = 0.0 // Safety against division by zero
	}

	// LOOP FILTER & NCO UPDATE
	state.CodePhase += state.CodeError * state.Gain
}
