package utils

import (
	"fmt"
	"time"
)

func NewTicker(name string, tickPerSecond int, profile bool, tolerancePercent float32) <-chan int64 {
	ticker := make(chan int64)
	go func(t chan int64) {
		acceptableProfileDurMax := time.Second + time.Second/100*time.Duration(tolerancePercent)
		acceptableProfileDurMin := time.Second - time.Second/100*time.Duration(tolerancePercent)
		fmt.Printf("tolerance %s: %v\n", name, acceptableProfileDurMax)
		tickDur := time.Second / time.Duration(tickPerSecond)
		fmt.Printf("expected Tick Duration for %s: %v\n", name, tickDur)
		tickDur -= tickDur / 100 * time.Duration(tolerancePercent)
		fmt.Printf("Compensated Tick Duration for %s: %v\n", name, tickDur)
		tickDurStart := time.Now()
		tickProfilerStart := tickDurStart
		tickProfilerCounter := 0
		for {
			now := time.Now()
			if now.Sub(tickDurStart) >= tickDur {
				tickDurStart = now
				t <- now.UnixNano()
				tickProfilerCounter++
				if tickProfilerCounter == tickPerSecond {
					tickProfilerCounter = 0
					if profile {
						profileDur := now.Sub(tickProfilerStart)
						if profileDur > acceptableProfileDurMax || profileDur < acceptableProfileDurMin {
							fmt.Printf("%s: %v, time: %v\n", name, time.Since(tickProfilerStart), now)
						}
						tickProfilerStart = now
					}
				}
			}
		}
	}(ticker)
	return ticker
}

func Idle() {
	time.Sleep(time.Microsecond)
}
