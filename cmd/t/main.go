package main

import (
	"fmt"
	"time"
)

func main() {
	const TIMESCALE float64 = 1e-9
	t1 := time.Now().UnixNano()
	time.Sleep(time.Millisecond * 1345)
	t2 := time.Now().UnixNano()
	passed := float64((t2 - t1)) * TIMESCALE
	fmt.Println(passed)
}
