package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/timeinterval"
)

func main() {
	ti := timeinterval.WithMinInterval(10, 50)
	counter := 0
	start := time.Now()
	for {
		if ti.IsTime() {
			counter++
			fmt.Printf("%d, %v\n", counter, time.Since(start))
		}
	}
}
