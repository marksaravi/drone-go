package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	hardware.InitHost()
	rxtxType := flag.String("t", "rx", "t")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	if *rxtxType == "rx" {
		runReceiver(ctx, &wg)
	} else {
		runTransmitter(ctx, &wg)
	}

	utils.WaitToAbortByENTER(cancel)
	wg.Wait()
}
