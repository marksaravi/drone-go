package plotter

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

const UDP_BUFFER_SIZE = 8192

type plotter struct {
	websocketConn     *websocket.Conn
	httpServerAddress string
	httpServer        *http.Server
	udpServerAddress  string
	udpConn           *net.UDPConn
	udpBuffer         []byte
	stopPlotter       context.CancelCauseFunc
}

func NewPlotter() *plotter {
	return &plotter{
		httpServerAddress: "localhost:3000",
		udpServerAddress:  "localhost:8000",
	}
}

func (p *plotter) StartPlotter() {
	log.SetFlags(log.Lmicroseconds)
	var ctx context.Context
	ctx, p.stopPlotter = context.WithCancelCause(context.Background())
	var wg sync.WaitGroup

	p.createUdpServer()
	p.createHttpServer()
	p.waitForInterrupt(ctx, &wg)
	p.stopByEnter()
	// p.startHttpServer(&wg)
	p.startUdpServer(&wg)
	log.Println("Plotter started...")
	wg.Wait()
	log.Printf("plotter stopped for: %v\n", context.Cause(ctx))
}

func (p *plotter) createUdpServer() {
	udpAddr, _ := net.ResolveUDPAddr("udp", p.udpServerAddress)
	var err error
	p.udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		p.stopPlotter(err)
		return
	}
	p.udpBuffer = make([]byte, UDP_BUFFER_SIZE)
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

func (p *plotter) startUdpServer(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting UDP Server ...")
		counter := 0
		lastPrint := time.Now()
		for {
			nBytes, _, err := p.udpConn.ReadFromUDP(p.udpBuffer)
			if err != nil && strings.Contains(err.Error(), "closed network connection") {
				return
			}
			if err == nil && nBytes > 0 {
				counter++
				if time.Since(lastPrint) >= time.Second {
					log.Printf("udp packet size: %d, packet per second: %d\n", nBytes, counter)
					counter = 0
					lastPrint = time.Now()
				}
			}
		}
	}()
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
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		defer p.udpConn.Close()
		log.Println("Stopping UDP Server ...")
		// defer p.shutdownHttpServer(wg)
		defer close(sigint)
		defer log.Printf("Stopping Plotter...")

		select {
		case <-sigint:
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
