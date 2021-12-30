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
