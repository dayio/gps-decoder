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

	sampleRate := 2000000.0
	chunkSize := int(sampleRate / 1000.0) // 2000 samples = 1ms

	// Create channels for our 32 satellites
	channels := make([]chan []complex128, 32)

	for prn := 1; prn <= 32; prn++ {
		// Define a buffer of 10 chunks
		channels[prn-1] = make(chan []complex128, 10)

		// Start the state machine goroutine for this PRN
		go gps.RunChannel(prn, channels[prn-1], sampleRate)
	}

	inputBuffer := make([]int8, chunkSize*2) // I and Q bytes

	for {
		err := signalSource.Read(inputBuffer)

		if err != nil {
			break // EOF or SDR error
		}

		outputComplex := dsp.ToComplex(inputBuffer)

		for i := 0; i < 32; i++ {
			channels[i] <- outputComplex
		}

		fmt.Print(".")
	}

	// Close channels to let goroutines exit cleanly
	for i := 0; i < 32; i++ {
		close(channels[i])
	}
}
