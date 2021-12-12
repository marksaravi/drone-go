package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
)

func WaitToAbortByENTER(cancel context.CancelFunc, wg *sync.WaitGroup) {
	log.Println("Press ENTER to abort")
	wg.Add(1)
	go func(cancel context.CancelFunc) {
		defer wg.Done()
		defer log.Println("Aborting by user ENTER")
		fmt.Scanln()
		cancel()
	}(cancel)
}
