package gps

import (
	"math"
	"math/cmplx"
)

type AcquireResult struct {
	PRN         int
	BestPhase   int
	BestDoppler float64
	SNR         float64
}

func Acquire(samples []complex128, prn int, sampleRate float64) AcquireResult {
	goldCode := GenerateGoldCode(prn)

	maxPower := 0.0 // Peak signak
	sumPower := 0.0 // Used to calculate the SNR
	count := 0      // Number of "samples"

	bestPhase := 0
	bestDoppler := 0.0

	// Doppler scanning from -10000Hz to +10000Hz with 500Hz increments
	for doppler := -10000.0; doppler <= 10000.0; doppler += 500 {

		// We slide the window across the samples
		for phase := range samples {

			var correlation complex128

			// We check our 1023-bit key against this specific window
			for i := range 1023 {

				// Calculate the signal index with the offset (modulo to wrap around)
				// This also maps the 1023 chips to the 2000 samples
				sigIdx := (phase + int(float64(i)*(sampleRate/1023000.0))) % len(samples)

				// Phase rotation to compensate for Doppler
				angle := 2 * math.Pi * doppler * (float64(i) / sampleRate)
				phasor := cmplx.Exp(complex(0, -angle))

				// Simple correlation, Signal * Code * Doppler compensation
				correlation += samples[sigIdx] * complex(goldCode[i], 0) * phasor
			}

			power := cmplx.Abs(correlation)

			// We add this power to the total noise baseline
			sumPower += power
			count++

			// Check if we find a new peak
			if power > maxPower {
				maxPower = power
				bestPhase = phase
				bestDoppler = doppler
			}
		}
	}

	// Calculate the average noise and the SNR
	avgNoise := sumPower / float64(count)
	snr := maxPower / avgNoise

	return AcquireResult{
		PRN:         prn,
		BestPhase:   bestPhase,
		BestDoppler: bestDoppler,
		SNR:         snr,
	}
}
