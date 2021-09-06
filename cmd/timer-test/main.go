package main

import (
	"fmt"
	"time"

	"github.com/MarkSaravi/drone-go/utils"
)

func main() {
	imu := utils.NewImu()
	const dataPerSecond int = int(3200 * 2.5)
	loopDur := time.Second / time.Duration(dataPerSecond)
	const SECONDS int = 5
	const TOTAL int = SECONDS * dataPerSecond
	var counter int = 0
	fmt.Println("Starting timer")
	start := time.Now()
	loopStart := start
	imu.ResetTime()
	for {
		now := time.Now()
		if now.Sub(loopStart) >= loopDur {
			imu.ReadRotations()
			counter++
			if counter == TOTAL {
				break
			}
			loopStart = now
		}
	}
	dur := time.Since(start).Seconds()
	seconds := float64(SECONDS)
	dev := (dur - seconds) * 100 / seconds
	fmt.Printf("Dur: %f, Dev: %%%f\n", dur, dev)
}
