package utils

import (
	"time"
)

func NewTicker(tickPerSecond int) <-chan int64 {
	ticker := make(chan int64)
	go func(t chan int64) {
		start := time.Now()
		dur := time.Second / time.Duration(tickPerSecond)
		for {
			now := time.Now()
			if now.Sub(start) >= dur {
				start = now
				t <- now.UnixNano()
			}
			time.Sleep(time.Microsecond)
		}
	}(ticker)
	return ticker
}
