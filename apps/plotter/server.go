package plotter

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

type plotter struct {
	sigint chan os.Signal
}

type PlotterSettings struct {
	UDPServerAddress  string
	HTTPServerAddress string
}

func NewPlotter(settings PlotterSettings) *plotter {
	return &plotter{
		sigint: make(chan os.Signal, 1),
	}
}

func (p *plotter) StartPlotter() {
	ctx, cancel := context.WithCancelCause(context.Background())
	var wg sync.WaitGroup
	p.waitForInterrupt(cancel, &wg)

	wg.Wait()
	fmt.Printf("plotter stopped for: %v", context.Cause(ctx))
}

func (p *plotter) waitForInterrupt(cancel context.CancelCauseFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		signal.Notify(p.sigint, os.Interrupt)
		<-p.sigint
		log.Printf("Stopping Plotter...")
		cancel(fmt.Errorf("signal interrupt"))
		close(p.sigint)
	}()
}
