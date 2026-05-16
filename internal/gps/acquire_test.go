package gps

import (
	"math"
	"testing"
)

func TestAcquireFFT(t *testing.T) {
	prn := 15
	sampleRate := 2000000.0
	expectedPhase := 500
	expectedDoppler := 1500.0

	goldCode := GenerateGoldCode(prn)

	// Manual stretch to simulate a 2 MHz signal
	rawSamples := make([]complex128, 2000)
	for i := 0; i < len(rawSamples); i++ {
		chipIdx := int(float64(i-expectedPhase)*1023000.0/sampleRate) % 1023
		if chipIdx < 0 {
			chipIdx += 1023
		}
		rawSamples[i] = complex(goldCode[chipIdx], 0)
	}

	// Doppler shift (we used -expectedDoppler so signal is shifted)
	simulatedSignal := ApplyDoppler(rawSamples, -expectedDoppler, sampleRate, 2000)

	// Acquisition
	result := AcquireFFT(simulatedSignal, prn, sampleRate)

	// Assert
	if result.BestPhase != expectedPhase {
		t.Errorf("Incorrect phase : Expected : %d,  got : %d", result.BestPhase, expectedPhase)
	}

	if math.Abs(result.BestDoppler-expectedDoppler) > 0.1 {
		t.Errorf("Incorrect Doppler : Expected : %f,  got : %f", result.BestDoppler, expectedDoppler)
	}
}
