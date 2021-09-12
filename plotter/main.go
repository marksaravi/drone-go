package plotter

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	SERVER_PORT int = 8081
	UDP_PORT    int = 6431
)

func Start() {

	dataChannel := make(chan float32, 10)

	startUDPReceiverServer(dataChannel)
	handler := createWebSocketHandler(dataChannel)
	http.Handle("/", http.FileServer(http.Dir("./plotter/static")))
	http.HandleFunc("/conn", handler)
	log.Println(fmt.Sprintf("Server is listening on port %d\n", SERVER_PORT))
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", SERVER_PORT), nil)
	log.Fatal(err)
}

func createWebSocketHandler(dataChannel chan float32) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Establishing connection")
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "Closing the connection")

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		ctx = c.CloseRead(ctx)

		for {
			select {
			case <-ctx.Done():
				c.Close(websocket.StatusNormalClosure, "")
				return
			case value := <-dataChannel:
				err = wsjson.Write(ctx, c, fmt.Sprintf("{\"value\": %f", value))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func startUDPReceiverServer(dataChannel chan float32) {
	go func() {
		var value float32 = 0
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: UDP_PORT,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			panic(err)
		}

		defer conn.Close()
		fmt.Printf("server listening %s\n", conn.LocalAddr().String())

		for {
			message := make([]byte, 20)
			// rlen, remote, err := conn.ReadFromUDP(message[:])
			_, _, err := conn.ReadFromUDP(message[:])
			if err != nil {
				panic(err)
			}

			// data := strings.TrimSpace(string(message[:rlen]))
			// fmt.Printf("received: %s from %s\n", data, remote)
			dataChannel <- value
			value += 0.1
		}
	}()
}
