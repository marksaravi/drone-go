package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	STOP_DATA int = iota
	STOP_CONNECTION
)

func startReceivers(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan int, connectionChan <-chan bool) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer log.Println("exiting receivers")
		var running bool = true
		var dataOpen bool = true
		var connectionOpen bool = true

		lastData := time.Now()
		lastConn := time.Now()

		for running || dataOpen || connectionOpen {
			select {
			case <-ctx.Done():
				running = false
			case value, ok := <-dataChan:
				if ok {
					if time.Since(lastData) >= time.Second/3 {
						log.Println("value: ", value)
						lastData = time.Now()
					}
				}
				if !ok && dataOpen {
					dataOpen = false
					log.Println("data channel is closed")
				}
			case connected, ok := <-connectionChan:
				if ok {
					if time.Since(lastConn) >= time.Second {
						log.Println("connection: ", connected)
						lastConn = time.Now()
					}
				}
				if !ok && connectionOpen {
					connectionOpen = false
					log.Println("connection channel is closed")
				}
			default:
			}
		}
	}()
}

func startTransmitters(ctx context.Context, wg *sync.WaitGroup, stopChan <-chan int) (<-chan int, <-chan bool) {
	wg.Add(1)
	dataChan := make(chan int)
	connectionChan := make(chan bool)
	go func() {
		defer wg.Done()
		defer log.Println("exiting transmitters")
		var running bool = true
		var stopOpen bool = true
		dataTimeout := time.Now()
		conTimeout := time.Now()
		for running || stopOpen {
			select {
			case <-ctx.Done():
				if running {
					log.Println("context DONE")
					running = false
				}
			case stop, ok := <-stopChan:
				if ok {
					switch stop {
					case STOP_DATA:
						log.Println("closing data channel...")
						close(dataChan)
						dataChan = nil
					case STOP_CONNECTION:
						log.Println("closing connection channel...")
						close(connectionChan)
						connectionChan = nil
					}
				}
				stopOpen = ok
			default:
			}
			if connectionChan != nil {
				if time.Since(conTimeout) > time.Millisecond*10 {
					conTimeout = time.Now()
					connectionChan <- true
				}
			}
			if dataChan != nil {
				if time.Since(dataTimeout) > time.Millisecond {
					dataTimeout = time.Now()
					dataChan <- rand.Intn(100)
				}
			}
		}
	}()
	return dataChan, connectionChan
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	var stopChan chan int = make(chan int)
	dataChan, connectionChan := startTransmitters(ctx, &wg, stopChan)
	startReceivers(ctx, &wg, dataChan, connectionChan)
	log.Println("waiting for command...")
	fmt.Scanln()
	stopChan <- STOP_DATA

	fmt.Scanln()
	stopChan <- STOP_CONNECTION
	close(stopChan)

	fmt.Scanln()
	log.Println("waiting for processes to stop...")
	cancel()
	wg.Wait()
	log.Println("program ended.")
}
