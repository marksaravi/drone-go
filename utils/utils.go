package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func WaitToAbortByENTER(cancel context.CancelFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("Press ENTER to abort")
	go func(cancel context.CancelFunc) {
		defer log.Println("Aborting by user ENTER")
		defer wg.Done()
		fmt.Scanln()
		cancel()
	}(cancel)
}

var lastWarning time.Time = time.Now()

func limitWarning(less bool, limit float64) {
	if time.Since(lastWarning) > time.Second/2 {
		lastWarning = time.Now()
		if less {
			log.Println("value is less than ", limit)
		} else {
			log.Println("value is more than ", limit)
		}
	}
}
func ApplyLimit(x, min, max float64, genwarning bool) float64 {
	if x < min {
		if genwarning {
			limitWarning(true, min)
		}
		return min
	}
	if x > max {
		if genwarning {
			limitWarning(false, max)
		}
		return max
	}
	return x
}
