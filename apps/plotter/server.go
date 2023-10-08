package plotter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"nhooyr.io/websocket"
)

type plotter struct {
	sigint            chan os.Signal
	websocketConn     *websocket.Conn
	httpServerAddress string
	httpServer        *http.Server
	stopPlotter       context.CancelCauseFunc
}

func NewPlotter() *plotter {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	return &plotter{
		sigint:            sigint,
		httpServerAddress: "localhost:3000",
	}
}

func (p *plotter) StartPlotter() {
	log.SetFlags(log.Lmicroseconds)
	var ctx context.Context
	ctx, p.stopPlotter = context.WithCancelCause(context.Background())
	var wg sync.WaitGroup

	p.createHttpServer()
	p.waitForInterrupt(ctx, &wg)
	p.stopByEnter()
	p.startHttpServer(&wg)
	log.Println("Plotter started...")
	wg.Wait()
	log.Printf("plotter stopped for: %v\n", context.Cause(ctx))
}

func (p *plotter) createHttpServer() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./apps/plotter/static")))
	mux.HandleFunc("/ws", p.createSocketHandler())
	p.httpServer = &http.Server{
		Addr:    p.httpServerAddress,
		Handler: mux,
	}
}

func (p *plotter) startHttpServer(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.Println("Starting HttpServer ...")
		if err := p.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			p.stopPlotter(err)
		}
	}()
}

func (p *plotter) shutdownHttpServer(wg *sync.WaitGroup) {
	log.Println("Shutting Down HttpServer...")
	p.httpServer.Shutdown(context.Background())
	wg.Done()
}

func (p *plotter) waitForInterrupt(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		defer log.Printf("Stopping Plotter...")
		defer p.shutdownHttpServer(wg)
		defer close(p.sigint)

		select {
		case <-p.sigint:
			p.stopPlotter(fmt.Errorf("signal interrupt"))
		case <-ctx.Done():
			return
		}
	}()
}

func (p *plotter) stopByEnter() {
	go func() {
		fmt.Scanln()
		p.stopPlotter(fmt.Errorf("enter"))
	}()
}

func (p *plotter) createSocketHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		p.websocketConn, err = websocket.Accept(w, r, nil)
		if err != nil {
			p.websocketConn = nil
			log.Println(err)
		} else {
			log.Println("Websocket Connection Accepted")
		}
	}
}
