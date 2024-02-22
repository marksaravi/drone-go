package utils

import "time"

type timeInterval struct {
	dataPerSecond int
	readInterval  time.Duration
	startTime     time.Time
	readCounter   int
}

func NewTimeInterval(dataPerSecond int) *timeInterval {
	return &timeInterval{
		dataPerSecond: dataPerSecond,
		readInterval:  time.Second / time.Duration(dataPerSecond),
		readCounter:   0,
		startTime:     time.Now(),
	}
}

func (m *timeInterval) IsTime() bool {
	n := int(time.Since(m.startTime) / m.readInterval)
	if n > m.readCounter {
		m.readCounter = n
		return true
	}
	return false
}
