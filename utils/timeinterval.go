package utils

import (
	"time"
)

type timeInterval struct {
	minIntervalPercent uint
	minInterval        time.Duration
	dataPerSecond      int
	readInterval       time.Duration
	startTime          time.Time
	readCounter        int
	lastRead           time.Time
}

func WithDataPerSecond(dataPerSecond int) *timeInterval {
	return WithMinInterval(dataPerSecond, 0)
}

func WithMinInterval(dataPerSecond int, minIntervalPercent uint) *timeInterval {
	if minIntervalPercent > 100 {
		minIntervalPercent = 100
	}
	readInterval := time.Second / time.Duration(dataPerSecond)
	minInterval := readInterval * time.Duration(minIntervalPercent) / 100
	return &timeInterval{
		minIntervalPercent: minIntervalPercent,
		minInterval:        minInterval,
		dataPerSecond:      dataPerSecond,
		readInterval:       readInterval,
		readCounter:        0,
		startTime:          time.Now(),
		lastRead:           time.Now(),
	}
}

func (m *timeInterval) IsTime() bool {
	if time.Since(m.lastRead) < m.minInterval {
		return false
	}
	n := int(time.Since(m.startTime) / m.readInterval)
	if n > m.readCounter {
		m.readCounter = n
		m.lastRead = time.Now()
		return true
	}
	return false
}
