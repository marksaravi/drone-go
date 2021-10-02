package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const LEVELS int = 5

func main() {
	log.SetFlags(log.Lmicroseconds)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go child(ctx, &wg, 0, nil)
	fmt.Scanln()
	log.Println("ENTER")
	cancel()
	wg.Wait()
}

func child(ctx context.Context, wg *sync.WaitGroup, level int, dc chan<- int) {
	childctx, cancel := context.WithCancel(ctx)
	data := make(chan int)
	ticker := time.NewTicker(time.Millisecond * time.Duration(317*level+317))
	defer wg.Done()
	defer log.Println("STOPPED ", level)
	defer close(data)
	defer ticker.Stop()

	if level < LEVELS {
		wg.Add(1)
		go child(childctx, wg, level+1, data)
	}
	log.Println("Starting ", level)

	for {
		select {
		case <-ctx.Done():
			cancel()
			log.Println("STOP FOR ", level)
			return
		case <-ticker.C:
			if dc != nil {
				dc <- level
			}

		case value := <-data:
			log.Println("from routine: ", value)
		}
	}
}
