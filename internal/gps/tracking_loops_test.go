package gps

import (
	"math"
	"testing"
)

func TestUpdatePLL(t *testing.T) {
	state := PLLState{
		DopplerFreq: 1500.0,
		Gain:        0.1,
	}

	// Simulate phase error (energy leaking into Q axis)
	I := 100.0
	Q := 20.0

	// Run loop to simulate continuous tracking over 20ms
	for ms := 0; ms < 20; ms++ {
		UpdatePLL(I, Q, &state)
	}

	// Doppler must increase to compensate for the positive phase error
	if state.DopplerFreq <= 1500.0 {
		t.Errorf("PLL failed to correct Doppler. Got: %f", state.DopplerFreq)
	}
}

func TestUpdateDLL(t *testing.T) {
	state := DLLState{
		CodePhase: 500.0,
		Gain:      0.05,
	}

	// Early > Late indicates local code is lagging behind the received signal
	earlyPower := 120.0
	latePower := 80.0

	UpdateDLL(earlyPower, latePower, &state)

	// Math check: Error = (120 - 80) / 200 = 0.2
	// Correction = 0.2 * 0.05 = 0.01
	expectedPhase := 500.01

	if math.Abs(state.CodePhase-expectedPhase) > 0.0001 {
		t.Errorf("DLL correction failed. Expected: %f, Got: %f", expectedPhase, state.CodePhase)
	}
}
