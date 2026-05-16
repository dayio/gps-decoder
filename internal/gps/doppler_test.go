package gps

import (
	"testing"
)

func TestApplyDoppler(t *testing.T) {
	// Fake signal (1, 2, 3)
	samples := []complex128{complex(1, 0), complex(2, 0), complex(3, 0)}
	padSize := 8
	doppler := 0.0 // No frequency shift
	sampleRate := 2000000.0

	result := ApplyDoppler(samples, doppler, sampleRate, padSize)

	// Size must comply with the padding size
	if len(result) != padSize {
		t.Errorf("Incorrect size. Expected : %d, got : %d", padSize, len(result))
	}

	// A zero-frequency shift should not change the signal
	for i, val := range samples {
		if result[i] != val {
			t.Errorf("Unexpected update at index %d. Expected : %v, got : %v", i, val, result[i])
		}
	}

	// The padding must be zero
	for i := len(samples); i < padSize; i++ {
		if result[i] != complex(0, 0) {
			t.Errorf("Padding failed at index %d. Got : %v", i, result[i])
		}
	}
}
