package utils

import (
	"context"
	"log"
	"runtime"
	"time"
)

func Idle() {
	runtime.Gosched()
}

func NewTicker(ctx context.Context, name string, tickPerSecond int, tolerancePercent float32) <-chan int64 {
	ticker := make(chan int64)
	go func() {
		defer log.Println("Closing Ticker (", name, ")")
		tickDur := time.Second / time.Duration(tickPerSecond)
		tickDur -= tickDur / 100 * time.Duration(tolerancePercent)
		tickDurStart := time.Now()
		for {
			select {
			case <-ctx.Done():
				close(ticker)
				return
			default:
				now := time.Now()
				if now.Sub(tickDurStart) >= tickDur {
					tickDurStart = now
					ticker <- now.UnixNano()
				}
				Idle()
			}
		}
	}()
	return ticker
}

type tickerProfiler struct {
	name          string
	tickPerSecond int
	tickCounter   int
	tickStart     time.Time
}

func NewTickerProfiler(name string, tickPerSecond int) *tickerProfiler {
	tp := tickerProfiler{
		name:          name,
		tickPerSecond: tickPerSecond,
	}
	tp.Reset()
	return &tp
}

func (tp *tickerProfiler) Count() {
	tp.tickCounter++
	if tp.tickCounter == tp.tickPerSecond {
		log.Printf("%s time for %d ticks: %v\n", tp.name, tp.tickPerSecond, time.Since(tp.tickStart))
		tp.Reset()
	}
}

func (tp *tickerProfiler) Reset() {
	tp.tickStart = time.Now()
	tp.tickCounter = 0
}
