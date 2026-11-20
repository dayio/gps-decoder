package signal

import (
	"fmt"
	"math"
	"time"

	"github.com/dayio/gps-decoder/internal/events"
	"github.com/dayio/gps-decoder/internal/gps"
)

type ChannelState int

const (
	StateAcquiring ChannelState = iota
	StateTracking
)

// acquireSemaphore limits concurrent FFT acquisitions to prevent CPU thrashing
var acquireSemaphore = make(chan struct{}, 4)

func RunChannel(prn int, sampleStream <-chan []complex128, sampleRate float64, bus *events.Bus, store *gps.NavStore) {
	state := StateAcquiring
	var trackState TrackingState

	navState := gps.NavDataState{}

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

					bus.Publish(events.Event{
						Timestamp: time.Now(),
						PRN:       prn,
						Type:      events.EventAcquired,
						Message:   fmt.Sprintf("Signal acquired (SNR: %.2f)", acquireResult.SNR),
					})

					initCodePhase := math.Mod(-float64(acquireResult.BestPhase)*1023_000.0/sampleRate, 1023.0)

					if initCodePhase < 0 {
						initCodePhase += 1023.0
					}

					trackState = TrackingState{
						PRN:          prn,
						CarrierPhase: 0.0,
						PLL: PLLState{
							DopplerFreq: acquireResult.BestDoppler,
							GainP:       0.25,
							GainI:       2.0,
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

					bus.Publish(events.Event{
						Timestamp: time.Now(),
						PRN:       prn,
						Type:      events.EventTrackingLost,
						Message:   "Loss of lock, falling back to acquisition",
					})

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
				gps.ProcessNavigationData(prn, trackResult.I, &navState, bus, store)
			}

			bus.Publish(events.Event{
				Timestamp: time.Now(),
				PRN:       prn,
				Type:      events.EventStateUpdate,
				Data: events.ChannelSnapshot{
					PRN:         prn,
					State:       "TRACKING",
					DopplerFreq: trackState.PLL.DopplerFreq,
					CodePhase:   trackState.DLL.CodePhase,
					PromptPower: trackResult.PromptPower,
					IsLocked:    trackResult.IsLocked,
				},
			})
		}
	}
}
