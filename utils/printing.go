package utils

import "time"

var startTimes map[string]time.Time = make(map[string]time.Time)

func PrintByInterval(id string, period time.Duration, printfn func()) {
	ts, exists := startTimes[id]
	if !exists {
		startTimes[id] = time.Now()
		return
	}
	if time.Since(ts) >= period {
		startTimes[id] = time.Now()
		printfn()
	}
}
