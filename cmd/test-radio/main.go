package main

import (
	"context"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	runTransmitter(ctx, &wg)
	utils.WaitToAbortByENTER(cancel)
	wg.Wait()
}
