package main

import (
	"fmt"
	"time"
)

func main() {
	sleepDur := time.Millisecond
	var numOfLoopsPerSecond int = int(time.Second / sleepDur)
	fmt.Printf("numOfLoopsPerSecond: %d\n", numOfLoopsPerSecond)
	start := time.Now()
	for i := 0; i < numOfLoopsPerSecond; i++ {
		time.Sleep(sleepDur)
	}
	actualDur := time.Since(start)
	expectedDur := sleepDur * time.Duration(numOfLoopsPerSecond)
	fmt.Printf("Expected Dur: %v, Actual Dur: %v, Ratio: %d\n", expectedDur, actualDur, actualDur/expectedDur)
}
