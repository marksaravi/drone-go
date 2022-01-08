package utils

import (
	"context"
	"fmt"
	"log"
)

func WaitToAbortByENTER(cancel context.CancelFunc) {
	log.Println("Press ENTER to abort")
	go func(cancel context.CancelFunc) {
		defer log.Println("Aborting by user ENTER")
		fmt.Scanln()
		cancel()
	}(cancel)
}

func ApplyLimit(x, min, max float64, genwarning bool) float64 {
	if x < min {
		if genwarning {
			log.Println("value is less than ", min)
		}
		return min
	}
	if x > max {
		if genwarning {
			log.Println("value is more than ", max)
		}
		return max
	}
	return x
}
