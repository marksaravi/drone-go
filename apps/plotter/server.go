package plotter

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	SERVER_PORT   int = 8081
	UDP_PORT      int = 6437
	IMU_DATA_SIZE int = 26
	TIME_SIZE     int = 8
)

func Start() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	dataChannel := make(chan string)
	startUDPReceiverServer(ctx, &wg, dataChannel)
	handler := createWebSocketHandler(ctx, dataChannel)
	http.Handle("/", http.FileServer(http.Dir("./apps/plotter/static")))
	http.HandleFunc("/conn", handler)
	var server = http.Server{
		Addr: fmt.Sprintf(":%d", SERVER_PORT),
	}

	go func() {
		log.Println(fmt.Sprintf("Server is listening on port %d\n", SERVER_PORT))
		if err := server.ListenAndServe(); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
	}()
	go func() {
		log.Println("Press ENTER to stop server")
		fmt.Scanln()
		fmt.Println("Stopping All Servers...")
		cancel()
	}()

	wg.Wait()
	fmt.Println("Stopping HTTP Server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}

func createWebSocketHandler(ctx context.Context, dataChannel chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Establishing connection")
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "Closing the connection")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("createWebSocketHandler done...")
				c.Close(websocket.StatusNormalClosure, "")
				return
			case json := <-dataChannel:
				err = wsjson.Write(ctx, c, json)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func startUDPReceiverServer(ctx context.Context, wg *sync.WaitGroup, dataChannel chan string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("End...")

		var value float32 = 0
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: UDP_PORT,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			panic(err)
		}

		log.Printf("UDP Server Listening %s\n", conn.LocalAddr().String())

		const bufferSize int = 16348
		data := make([]byte, bufferSize)

		running := true

		for running {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping UDP Server...")
				conn.Close()
				running = false
			default:
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				nBytes, _, err := conn.ReadFromUDP(data)
				if err == nil && nBytes > 0 {
					dataPerPacket := dataPerPacket(data[0:2])
					fmt.Println(dataPerPacket)
					jsonData := extractPackets(data[2:], dataPerPacket)
					// fmt.Println(jsonData)
					dataChannel <- jsonData

				}
				value += 0.1
			}
		}
	}()
}

func extractPackets(data []byte, dataPerPacket int) string {
	var jsonData string = "["
	var comma string = ""
	for i := 0; i < 256; i++ {
		imudata := extractImuRotations(data[i*26 : (i+1)*26])
		jsonData += comma + imudata
		comma = ","
	}
	return jsonData + "]"
}

func extractImuRotations(data []byte) string {
	t := binary.LittleEndian.Uint64(data[0:8])
	a := extractRotations(data[8:14])
	g := extractRotations(data[14:20])
	r := extractRotations(data[20:26])
	return fmt.Sprintf("{\"a\":%s,\"g\":%s,\"r\":%s,\"t\":%d}", a, g, r, t)
}

func extractRotations(data []byte) string {
	rot := make([]float64, 3)
	for i := 0; i < 3; i++ {
		rot[i] = rpy(data[8+i*2 : 10+i*2])
	}
	return fmt.Sprintf("{\"roll\":%0.2f,\"pitch\":%0.2f,\"yaw\":%0.2f}", rot[0], rot[1], rot[2])
}

func dataPerPacket(data []byte) int {
	return int(binary.LittleEndian.Uint16(data))
}

func rpy(data []byte) float64 {
	i := binary.LittleEndian.Uint16(data)
	i -= 16000
	return float64(i) / 10
}
