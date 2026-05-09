package main

import (
	"fmt"
	"log"

	"github.com/dayio/gps-decoder/internal/dsp"
	"github.com/dayio/gps-decoder/internal/gps"
	"github.com/dayio/gps-decoder/internal/source"
)

func main() {

	var signalSource source.IQSource
	var err error

	signalSource, err = source.ReadFile("./data/gpssim.bin")

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	defer func() {
		if err := signalSource.Close(); err != nil {
			log.Printf("Error when closing : %v\n", err)
		}
	}()

	inputBuffer := make([]int8, 4000)

	for {
		err := signalSource.Read(inputBuffer)

		if err != nil {
			break // EOF or SDR error
		}

		outputComplex := dsp.ToComplex(inputBuffer)

		sampleRate := 2000000.0

		for prn := 1; prn <= 32; prn++ {

			bestPhase, bestDoppler, snr := gps.Acquire(outputComplex, prn, sampleRate)

			if snr > 3.0 {
				fmt.Printf("Satellite PRN %02d found ! Phase: %4d | Doppler: %5.2f | SNR: %5.2f\n", prn, bestPhase, bestDoppler, snr)
			}
		}
	}
}
