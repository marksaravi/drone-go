package drone

import (
	"context"
	"log"
	"sync"
)

type drone struct {
}

func NewDrone() *drone {
	return &drone{}
}

func (d *drone) Start(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("drone started")
	defer log.Println("drone stopped")
	d.controller(ctx, wg)
}

func (d *drone) controller(ctx context.Context, wg *sync.WaitGroup) {
	running := true
	for running {
		select {
		case <-ctx.Done():
			log.Println("STOP SIGNAL RECEIVED")
			running = false
		default:
		}
	}
}
