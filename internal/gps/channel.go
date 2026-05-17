package gps

import (
	"fmt"
)

type ChannelState int

const (
	StateAcquiring ChannelState = iota
	StateTracking
)

func RunChannel(prn int, sampleStream <-chan []complex128, sampleRate float64) {
	state := StateAcquiring
	var trackState TrackingState

	navState := NavDataState{}

	// Sleep penality
	skipCounter := 0

	for chunk := range sampleStream {

		if skipCounter > 0 {
			skipCounter--
			continue
		}

		switch state {

		case StateAcquiring:
			acquireResult := AcquireFFT(chunk, prn, sampleRate)

			// Threshold to switch to tracking
			if acquireResult.SNR > 6.0 {
				fmt.Printf("\n[PRN %02d] Signal acquired - SNR: %.2f. Switching to TRACKING\n", prn, acquireResult.SNR)

				trackState = TrackingState{
					PRN:          prn,
					CarrierPhase: 0.0,
					PLL: PLLState{
						DopplerFreq: acquireResult.BestDoppler,
						Gain:        0.25,
					},
					DLL: DLLState{
						CodePhase: float64(acquireResult.BestPhase),
						Gain:      0.01,
					},
				}
				state = StateTracking
			} else {
				skipCounter = 100 // Allow some time to recover from the noise
			}

		case StateTracking:
			trackResult := TrackChunk(chunk, &trackState, sampleRate)

			// Loss of Lock detection
			if trackResult.PromptPower < 50.0 {
				fmt.Printf("\n[PRN %02d] Tracking lost - Fallback to ACQUISITION\n", prn)
				state = StateAcquiring

				skipCounter = 500 // Sleep for 500ms let the satellite signal recover
				continue
			}

			// Data Extraction
			if trackResult.IsLocked {
				ProcessNavigationData(prn, trackResult.I, &navState)
			}
		}
	}
}
