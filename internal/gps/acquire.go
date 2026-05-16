package gps

import (
	"math/cmplx"

	"github.com/madelynnblue/go-dsp/fft"
)

type AcquireResult struct {
	PRN         int
	BestPhase   int
	BestDoppler float64
	SNR         float64
}

// AcquireFFT remplace l'ancienne méthode brute-force par le théorème de la corrélation circulaire
func AcquireFFT(samples []complex128, prn int, sampleRate float64) AcquireResult {
	goldCode := GenerateGoldCode(prn)

	N := 2048

	paddedCode := make([]complex128, N)

	for i := 0; i < len(samples); i++ {
		// stretch the signal to fit the 1023-bit code
		chipIdx := int(float64(i)*1023_000.0/sampleRate) % 1023
		paddedCode[i] = complex(goldCode[chipIdx], 0)
	}

	codeFFT := fft.FFT(paddedCode)

	for i := range codeFFT {
		codeFFT[i] = cmplx.Conj(codeFFT[i])
	}

	maxPower := 0.0
	sumPower := 0.0
	count := 0
	bestPhase := 0
	bestDoppler := 0.0

	// Doppler scanning
	for doppler := -10000.0; doppler <= 10000.0; doppler += 500 {

		// Rotation and padding
		rotatedSamples := ApplyDoppler(samples, doppler, sampleRate, N)

		// Go to frequency domain
		signalFFT := fft.FFT(rotatedSamples)

		// Circular correlation (Multiplication)
		for i := range signalFFT {
			signalFFT[i] = signalFFT[i] * codeFFT[i]
		}

		// Back to temporal domain
		correlationScores := fft.IFFT(signalFFT)

		// Peak correlation analysis without padding values
		for i := 0; i < len(samples); i++ {
			power := cmplx.Abs(correlationScores[i])

			sumPower += power
			count++

			if power > maxPower {
				maxPower = power
				bestPhase = i
				bestDoppler = doppler
			}
		}
	}

	avgNoise := sumPower / float64(count)
	snr := maxPower / avgNoise

	return AcquireResult{
		PRN:         prn,
		BestPhase:   bestPhase,
		BestDoppler: bestDoppler,
		SNR:         snr,
	}
}
