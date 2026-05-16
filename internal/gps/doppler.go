package gps

import (
	"math"
	"math/cmplx"
)

func ApplyDoppler(samples []complex128, doppler float64, sampleRate float64, padSize int) []complex128 {
	rotated := make([]complex128, padSize)

	for i := range samples {
		angle := 2 * math.Pi * doppler * (float64(i) / sampleRate)
		rotated[i] = samples[i] * cmplx.Exp(complex(0, -angle))
	}

	return rotated
}
