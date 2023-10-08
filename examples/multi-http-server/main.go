package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type plotter struct {
	sigint      chan os.Signal
	httpServers map[int]*http.Server
}

func MultiServer() *plotter {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	return &plotter{
		sigint:      sigint,
		httpServers: make(map[int]*http.Server),
	}
}

func (p *plotter) StartMultiServer(numberOfServers int) {
	log.SetFlags(log.Lmicroseconds)
	ctx, cancel := context.WithCancelCause(context.Background())
	var wg sync.WaitGroup
	for serverIndex := 0; serverIndex < numberOfServers; serverIndex++ {
		p.createHttpServer(serverIndex)
	}
	p.waitForInterrupt(ctx, cancel, &wg)
	p.stopByEnter(cancel)
	p.startHttpServers(&wg, cancel)
	log.Println("Plotter started...")
	wg.Wait()
	log.Printf("plotter stopped for: %v\n", context.Cause(ctx))
}

func (p *plotter) createHttpServer(index int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", createFuncHandler(fmt.Sprintf("Server#%d", index)))
	p.httpServers[index] = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", 3000+index),
		Handler: mux,
	}
}

func (p *plotter) startHttpServers(wg *sync.WaitGroup, cancel context.CancelCauseFunc) {
	for index, httpServer := range p.httpServers {
		wg.Add(1)
		go func(i int, server *http.Server) {
			log.Printf("Starting HttpServer #%d...\n", i)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				cancel(err)
			}
		}(index, httpServer)
	}
}

func (p *plotter) shutdownHttpServers(wg *sync.WaitGroup) {
	for index, server := range p.httpServers {
		log.Printf("Shutting Down HttpServer #%d...\n", index)
		server.Shutdown(context.Background())
		wg.Done()
	}

}

func (p *plotter) waitForInterrupt(ctx context.Context, cancel context.CancelCauseFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer log.Printf("Stopping Plotter...")
		defer p.shutdownHttpServers(wg)
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

func createFuncHandler(id string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(fmt.Sprintf("from server %s: %d", id, time.Now().UnixMilli())))
	}
}

func main() {
	plotter := MultiServer()
	plotter.StartMultiServer(4)
}
