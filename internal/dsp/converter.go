package dsp

func ToComplex(buffer []int8) []complex128 {

	samples := make([]complex128, len(buffer)/2)

	for i := range samples {
		samples[i] = complex(float64(buffer[2*i]), float64(buffer[2*i+1]))
	}

	return samples
}
