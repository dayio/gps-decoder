package gps

import (
	"time"
)

type SubframeHeader struct {
	SubframeID int       // 1 à 5 (identifief in the HOW word)
	TOWCount   uint32    // raw counter (17 bits)
	TOW        float64   // Time of week in seconds (TOWCount * 6s)
	Timestamp  time.Time // Local time of reception
}

type Satellite struct {
	PRN         int
	LastHeader  SubframeHeader // Last validated TLM/HOW
	LastUpdated time.Time
}

type NavStore struct {
	Satellites map[int]*Satellite
}

func NewNavStore() *NavStore {
	return &NavStore{
		Satellites: make(map[int]*Satellite),
	}
}
