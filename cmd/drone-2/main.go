package main

import (
	"context"
	"log"
	"sync"

	dronepackage "github.com/marksaravi/drone-go/apps/drone"
	"github.com/marksaravi/drone-go/utils"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	utils.WaitToAbortByESC(cancel)
	drone := dronepackage.NewDrone()
	drone.Start(ctx, &wg)
}
