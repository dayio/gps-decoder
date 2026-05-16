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

	for chunk := range sampleStream {
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
						Gain:      0.05,
					},
				}
				state = StateTracking
			}

		case StateTracking:
			trackResult := TrackChunk(chunk, &trackState, sampleRate)

			// Loss of Lock detection
			if trackResult.PromptPower < 50.0 {
				fmt.Printf("\n[PRN %02d] Tracking lost - Fallback to ACQUISITION\n", prn)
				state = StateAcquiring
				continue
			}

			// Data Extraction
			if trackResult.IsLocked {
				// TODO: We will build this in the next article !
			}
		}
	}
}
