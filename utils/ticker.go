package utils

import (
	"log"
	"runtime"
	"time"
)

func NewTicker(name string, tickPerSecond int, tolerancePercent float32, enableProfiling bool) <-chan int64 {
	ticker := make(chan int64)
	go func(t chan int64) {
		acceptableProfileDurMax := time.Second + time.Second/100*time.Duration(tolerancePercent)
		acceptableProfileDurMin := time.Second - time.Second/100*time.Duration(tolerancePercent)
		log.Printf("tolerance %s: %v\n", name, acceptableProfileDurMax)
		tickDur := time.Second / time.Duration(tickPerSecond)
		log.Printf("expected Tick Duration for %s: %v\n", name, tickDur)
		tickDur -= tickDur / 100 * time.Duration(tolerancePercent)
		log.Printf("Compensated Tick Duration for %s: %v\n", name, tickDur)
		tickDurStart := time.Now()
		tickProfilerStart := tickDurStart
		tickProfilerCounter := 0
		for {
			now := time.Now()
			if now.Sub(tickDurStart) >= tickDur {
				tickDurStart = now
				t <- now.UnixNano()
				if enableProfiling {
					tickProfilerCounter++
					if tickProfilerCounter == tickPerSecond {
						tickProfilerCounter = 0
						profileDur := now.Sub(tickProfilerStart)
						if profileDur > acceptableProfileDurMax || profileDur < acceptableProfileDurMin {
							log.Printf("%s: %v\n", name, time.Since(tickProfilerStart))
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
	runtime.Gosched()
}
