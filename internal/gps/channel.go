package gps

import (
	"fmt"
	"math"
)

type ChannelState int

const (
	StateAcquiring ChannelState = iota
	StateTracking
)

// acquireSemaphore limits concurrent FFT acquisitions to prevent CPU thrashing
var acquireSemaphore = make(chan struct{}, 4)

func RunChannel(prn int, sampleStream <-chan []complex128, sampleRate float64) {
	state := StateAcquiring
	var trackState TrackingState

	navState := NavDataState{}

	// Sleep penality
	skipCounter := 0
	lossOfLockCount := 0

	for chunk := range sampleStream {

		if skipCounter > 0 {
			skipCounter--
			continue
		}

		switch state {

		case StateAcquiring:
			select {
			case acquireSemaphore <- struct{}{}:
				acquireResult := AcquireFFT(chunk, prn, sampleRate)
				<-acquireSemaphore

				// Threshold to switch to tracking
				if acquireResult.SNR > 8.0 {
					fmt.Printf("[PRN %02d] Signal acquired - SNR: %.2f. Switching to TRACKING\n", prn, acquireResult.SNR)

					initCodePhase := math.Mod(-float64(acquireResult.BestPhase)*1023_000.0/sampleRate, 1023.0)

					if initCodePhase < 0 {
						initCodePhase += 1023.0
					}

					trackState = TrackingState{
						PRN:          prn,
						CarrierPhase: 0.0,
						PLL: PLLState{
							DopplerFreq: acquireResult.BestDoppler,
							Gain:        0.25,
						},
						DLL: DLLState{
							CodePhase: initCodePhase,
							Gain:      0.01,
						},
					}
					state = StateTracking
					lossOfLockCount = 0
				} else {
					skipCounter = 100 // Allow some time to recover from the noise
				}
			default:
				skipCounter = 20 // Pool is busy, wait for the next cycle
			}

		case StateTracking:
			trackResult := TrackChunk(chunk, &trackState, sampleRate)

			// Loss of Lock detection
			if trackResult.PromptPower < 20.0 {
				lossOfLockCount++
				if lossOfLockCount > 100 {
					fmt.Printf("[PRN %02d] Tracking lost - Fallback to ACQUISITION\n", prn)
					state = StateAcquiring
					lossOfLockCount = 0
					skipCounter = 300 // Sleep for 300ms let the satellite signal recover
					continue
				}
			} else {
				lossOfLockCount = 0
			}

			// Data Extraction
			if trackResult.IsLocked {
				ProcessNavigationData(prn, trackResult.I, &navState)
			}
		}
	}
}
