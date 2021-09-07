package main

import (
	"fmt"
	"time"
)

func main() {
	N := 100000
	dur := time.Second / time.Duration(N)
	fmt.Println(dur)
	ticker := createTicker(N)
	counter := 0
	start := time.Now()
	for range ticker {
		counter++
		if counter == N {
			fmt.Println(time.Since(start))
			start = time.Now()
			counter = 0
		}
	}
}

func createTicker(tickPerSecond int) chan int64 {
	ticker := make(chan int64)
	go func(t chan int64) {
		start := time.Now()
		dur := time.Second / time.Duration(tickPerSecond)
		fmt.Println("Duration: ", dur)
		for {
			now := time.Now()
			if now.Sub(start) >= dur {
				start = now
				t <- now.UnixNano()
			}
		}
	}(ticker)
	return ticker
}
