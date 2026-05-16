package gps

import (
	"math"
	"math/cmplx"
)

type TrackStepResult struct {
	PromptPower float64
	I           float64
	IsLocked    bool
}

type TrackingState struct {
	PRN          int
	CarrierPhase float64 // Accumulated carrier phase
	PLL          PLLState
	DLL          DLLState
}

func TrackChunk(chunk []complex128, state *TrackingState, sampleRate float64) TrackStepResult {
	// Spacing for Early and Late correlators (0.5 chip is standard)
	codeSpacing := 0.5

	// CORRELATE (Integrate & Dump)
	// We test the local gold code at three distinct positions
	early := Correlate(chunk, state.PRN, state.DLL.CodePhase-codeSpacing, state.PLL.DopplerFreq, state.CarrierPhase, sampleRate)
	prompt := Correlate(chunk, state.PRN, state.DLL.CodePhase, state.PLL.DopplerFreq, state.CarrierPhase, sampleRate)
	late := Correlate(chunk, state.PRN, state.DLL.CodePhase+codeSpacing, state.PLL.DopplerFreq, state.CarrierPhase, sampleRate)

	// DLL UPDATE (Fix Code Delay)
	earlyPower := cmplx.Abs(early)
	latePower := cmplx.Abs(late)
	UpdateDLL(earlyPower, latePower, &state.DLL)

	// PLL UPDATE (Fix Doppler Frequency)
	I := real(prompt)
	Q := imag(prompt)
	UpdatePLL(I, Q, &state.PLL)

	// CARRIER PHASE ADVANCE
	// We move the carrier phase forward by 1ms for the next chunk
	state.CarrierPhase += 2.0 * math.Pi * state.PLL.DopplerFreq * 0.001
	state.CarrierPhase = math.Mod(state.CarrierPhase, 2*math.Pi) // Keep it within 0 to 2PI as the phase rotates

	// RETURN METRICS
	// If the phase error is small enough, the PLL is considered as "locked"
	isLocked := math.Abs(state.PLL.PhaseError) < 0.2

	return TrackStepResult{
		PromptPower: cmplx.Abs(prompt),
		I:           I,
		IsLocked:    isLocked,
	}
}

func Correlate(chunk []complex128, prn int, codePhase, doppler, initialPhase, sampleRate float64) complex128 {

	goldCode := GenerateGoldCode(prn)

	var correlation complex128

	for i, sample := range chunk {

		// Calculate the Gold Code for this specific sample
		chipIdx := int(codePhase+float64(i)*(1023_000.0/sampleRate)) % 1023

		if chipIdx < 0 {
			chipIdx += 1023
		}

		// Generate the local carrier wave to cancel the Doppler
		angle := initialPhase + 2.0*math.Pi*doppler*(float64(i)/sampleRate)

		// Negative angle to untwist the signal
		phasor := cmplx.Exp(complex(0, -angle))

		// Multiply incoming signal by local code and local carrier
		correlation += sample * complex(goldCode[chipIdx], 0) * phasor
	}

	return correlation
}
