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

func NewTicker(ctx context.Context, name string, tickPerSecond int) <-chan int64 {
	ticker := make(chan int64)
	go func() {
		defer log.Println("Ticker stopped (", name, ")")
		tickDur := int64(time.Second / time.Duration(tickPerSecond))
		if tickPerSecond >= 1000 {
			tickDur -= tickDur / 20
		}
		tickStart := time.Now().UnixNano()
		for ticker != nil {
			now := time.Now().UnixNano()
			if now-tickStart >= tickDur {
				tickStart = now
				ticker <- tickStart
			}
			select {
			case <-ctx.Done():
				close(ticker)
				ticker = nil
			default:
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
