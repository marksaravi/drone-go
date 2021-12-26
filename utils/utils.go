package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func WaitToAbortByENTER(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	log.Println("Press ENTER to abort")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(250 * time.Millisecond)
			}
		}
	}()
	go func(cancel context.CancelFunc) {
		defer log.Println("Aborting by user ENTER")
		fmt.Scanln()
		cancel()
	}(cancel)
}
