package utils

import (
	"fmt"
	"time"
)

func NewTicker(name string, tickPerSecond int, profile bool) <-chan int64 {
	ticker := make(chan int64)
	go func(t chan int64) {
		acceptableProfileDur := time.Second + time.Second/10
		tickDur := time.Second / time.Duration(tickPerSecond)
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
						if profileDur > acceptableProfileDur {
							fmt.Printf("%s: %v\n", name, time.Since(tickProfilerStart))
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
