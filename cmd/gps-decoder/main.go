package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./data/gpssim.bin")

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	buffer := make([]int8, 4000)

	err = binary.Read(file, binary.LittleEndian, &buffer)

	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	for i := 0; i < 10; i += 2 {
		signalI := buffer[i]
		signalQ := buffer[i+1]
		fmt.Printf("sample %d -> I: %4d, Q: %4d\n", (i/2)+1, signalI, signalQ)
	}
}
