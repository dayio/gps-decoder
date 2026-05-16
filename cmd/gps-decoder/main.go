package main

import (
	"fmt"
	"log"
	"sync"
	"time"

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
		now := time.Now()

		err := signalSource.Read(inputBuffer)

		if err != nil {
			break // EOF or SDR error
		}

		outputComplex := dsp.ToComplex(inputBuffer)

		sampleRate := 2000000.0

		var wg sync.WaitGroup

		acquireCh := make(chan gps.AcquireResult, 32)

		for prn := 1; prn <= 32; prn++ {
			wg.Add(1)

			go func(prn int) {
				defer wg.Done()
				acquireCh <- gps.AcquireFFT(outputComplex, prn, sampleRate)
			}(prn)
		}

		go func() {
			wg.Wait()
			close(acquireCh)
		}()

		for acquireResult := range acquireCh {
			if acquireResult.SNR > 8.0 {
				fmt.Printf("Satellite PRN %02d found ! Phase: %4d | Doppler: %5.2f | SNR: %5.2f\n",
					acquireResult.PRN,
					acquireResult.BestPhase,
					acquireResult.BestDoppler,
					acquireResult.SNR,
				)
			}
		}

		log.Printf("Processing took %v\n", time.Since(now))
	}
}
