package events

import "time"

type EventType string

const (
	EventAcquired      EventType = "ACQUIRED"
	EventTrackingLost  EventType = "TRACKING_LOST"
	EventPreambleFound EventType = "PREAMBLE_FOUND"
	EventStateUpdate   EventType = "STATE_UPDATE" // Noisy but useful for debugging
	EventWordDecoded   EventType = "WORD_DECODED"
)

type PreambleType int

const (
	PreambleNormal PreambleType = iota
	PreambleInverted
)

// Event représente un fait ponctuel survenu dans le décodeur
type Event struct {
	Timestamp time.Time
	PRN       int
	Type      EventType
	Message   string
	Data      any // optional payload (eg: NavDataState)
}

// ChannelSnapshot is a view of a satellite at a given time
type ChannelSnapshot struct {
	PRN          int
	State        string // "IDLE", "ACQUIRING", "TRACKING"
	SNR          float64
	DopplerFreq  float64
	CodePhase    float64
	PromptPower  float64
	IsLocked     bool
	BitBuffer    uint8
	PreambleSeen int
}
