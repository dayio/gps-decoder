package main

import (
	"fmt"
	"log"

	"github.com/dayio/gps-decoder/internal/source"
)

func main() {

	var signalSource source.IQSource
	var err error

	signalSource, err = source.ReadFile("./data/gpssim.bin")

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	defer signalSource.Close()

	buffer := make([]int8, 4000)

	for {
		err := signalSource.Read(buffer)

		if err != nil {
			break // EOF or SDR error
		}

		fmt.Println("sample", buffer)
	}
}
