package main

import (
	"fmt"
	"log"
	"sync"

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

	sampleRate := 2000000.0
	chunkSize := int(sampleRate / 1000.0) // 2000 samples = 1ms

	// Create channels for our 32 satellites
	channels := make([]chan []complex128, 32)
	var wg sync.WaitGroup

	for prn := 1; prn <= 32; prn++ {
		channels[prn-1] = make(chan []complex128, 1000)
		wg.Add(1)

		go func(p int, ch chan []complex128) {
			defer wg.Done()
			gps.RunChannel(p, ch, sampleRate)
		}(prn, channels[prn-1])
	}

	inputBuffer := make([]int8, chunkSize*2) // I and Q bytes
	processedMs := 0

	for {
		err := signalSource.Read(inputBuffer)

		if err != nil {
			break // EOF or SDR error
		}

		processedMs++

		if processedMs%5000 == 0 {
			fmt.Printf("\n--- PROGRESS: %.1f seconds processed ---\n", float64(processedMs)/1000.0)
		}

		outputComplex := dsp.ToComplex(inputBuffer)

		for i := range 32 {
			channels[i] <- outputComplex
		}
	}

	// Close channels to let goroutines exit cleanly
	for i := range 32 {
		close(channels[i])
	}

	wg.Wait()
}
