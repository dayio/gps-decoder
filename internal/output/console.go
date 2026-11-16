package output

import (
	"fmt"

	"github.com/dayio/gps-decoder/internal/events"
)

func RunConsoleSink(eventsCh <-chan events.Event) {
	for e := range eventsCh {
		if e.Type != events.EventStateUpdate {
			// fmt.Printf("[%s] [PRN %02d] %s: %s\n", e.Timestamp.Format("15:04:05.000"), e.PRN, e.Type, e.Message)
			fmt.Printf("[PRN %02d] %s: %s\n", e.PRN, e.Type, e.Message)
		}
	}
}
