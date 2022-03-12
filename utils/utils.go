package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

var intervals map[string]time.Time = make(map[string]time.Time)

func WaitToAbortByENTER(cancel context.CancelFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Press ENTER to abort")
	go func(cancel context.CancelFunc) {
		defer log.Println("Aborting by user ENTER")
		defer wg.Done()
		fmt.Scanln()
		cancel()
	}(cancel)
}

func ApplyLimits(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func PrintIntervally(msg string, id string, interval time.Duration, useLog bool) {
	log.SetFlags(log.Lmicroseconds)
	now := time.Now()
	last, ok := intervals[id]
	if !ok {
		last = now
		intervals[id] = now
	}
	if time.Since(last) >= interval {
		intervals[id] = now
		if useLog {
			log.Print(msg)
		} else {
			fmt.Print(msg)
		}
	}
}
