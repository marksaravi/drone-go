package plotter

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/marksaravi/drone-go/utils"
)

const (
	SERVER_PORT   int = 3000
	UDP_PORT      int = 4000
	IMU_DATA_SIZE int = 26
	TIME_SIZE     int = 8
)

func Start() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	dataChannel := make(chan string)
	startUDPReceiverServer(ctx, &wg, dataChannel)
	// http.Handle("/", http.FileServer(http.Dir("./apps/plotter/static")))
	// handler := createWebSocketHandler(ctx, dataChannel)

	// http.HandleFunc("/conn", handler)
	// var server = http.Server{
	// 	Addr: fmt.Sprintf(":%d", SERVER_PORT),
	// }

	go func() {
		log.Println("Press ENTER to stop server")
		fmt.Scanln()
		fmt.Println("Stopping All Servers...")
		cancel()
	}()

	log.Println(fmt.Sprintf("Server is listening on port %d\n", SERVER_PORT))
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Printf("HTTP server Shutdown: %v", err)
	// }
	wg.Wait()
	// fmt.Println("Stopping HTTP Server...")
	// if err := server.Shutdown(ctx); err != nil {
	// 	log.Fatal(err)
	// }

}

// func createWebSocketHandler(ctx context.Context, dataChannel chan string) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Establishing connection")
// 		c, err := websocket.Accept(w, r, nil)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				c.Close(websocket.StatusInternalError, "closing the websocket connection...")
// 				fmt.Println("Stopping WebSocketHandler...")
// 				return
// 			case json := <-dataChannel:
// 				err = wsjson.Write(ctx, c, json)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		}
// 	}
// }

func startUDPReceiverServer(ctx context.Context, wg *sync.WaitGroup, dataChannel chan string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("End...")

		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: UDP_PORT,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if err != nil {
			return
		}

		log.Printf("UDP Server Listening %s\n", conn.LocalAddr().String())

		const bufferSize int = 8192
		data := make([]byte, bufferSize)

		for {
			select {
			case <-ctx.Done():
				conn.Close()
				fmt.Println("Stopping UDP Server...")
				return
			default:
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				nBytes, _, err := conn.ReadFromUDP(data)
				if err == nil && nBytes > 0 {
					packetSize := int(utils.DeSerializeInt(data[0:2]))
					data = data[0:packetSize]
					dataPerPacket := int(utils.DeSerializeInt(data[2:4]))
					dataLen := int(utils.DeSerializeInt(data[4:6]))
					fmt.Println(packetSize, dataPerPacket, dataLen)
					jsonData := extractPackets(data[6:], dataLen, dataPerPacket)
					fmt.Println(jsonData[0:32])

					// dataChannel <- jsonData
				}
			}
		}
	}()
}

func extractPackets(data []byte, dataLen, dataPerPacket int) string {
	var jsonData string = "["
	var comma string = ""
	for i := 0; i < dataPerPacket; i++ {
		imudata := extractImuRotations(data[i*dataLen : (i+1)*dataLen])
		jsonData += comma + imudata
		comma = ","
	}
	return jsonData + "]"
}

func extractImuRotations(data []byte) string {
	t := utils.DeSerializeDuration(data[0:4])
	a := extractRotations(data[4:10])
	g := extractRotations(data[10:16])
	r := extractRotations(data[16:22])
	return fmt.Sprintf("{\"a\":%s,\"g\":%s,\"r\":%s,\"t\":%d}", a, g, r, t)
}

func extractRotations(data []byte) string {
	rot := make([]float64, 3)
	for i := 0; i < 3; i++ {
		rot[i] = utils.DeSerializeFloat64(data[i*2 : (i+1)*2])
	}
	return fmt.Sprintf("{\"roll\":%0.2f,\"pitch\":%0.2f,\"yaw\":%0.2f}", rot[0], rot[1], rot[2])
}
