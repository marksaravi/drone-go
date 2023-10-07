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
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	return &plotter{
		sigint: sigint,
	}
}

func (p *plotter) StartPlotter() {
	log.SetFlags(log.Lmicroseconds)
	ctx, cancel := context.WithCancelCause(context.Background())
	var wg sync.WaitGroup
	p.waitForInterrupt(ctx, cancel, &wg)
	p.stopByEnter(cancel)
	log.Println("Plotter started...")
	wg.Wait()
	log.Printf("plotter stopped for: %v\n", context.Cause(ctx))
}

func (p *plotter) waitForInterrupt(ctx context.Context, cancel context.CancelCauseFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer log.Printf("Stopping Plotter...")
		defer close(p.sigint)
		defer wg.Done()

		select {
		case <-p.sigint:
			cancel(fmt.Errorf("signal interrupt"))
		case <-ctx.Done():
			return
		}
	}()
}

func (p *plotter) stopByEnter(cancel context.CancelCauseFunc) {
	go func() {
		fmt.Scanln()
		cancel(fmt.Errorf("enter"))
	}()
}
