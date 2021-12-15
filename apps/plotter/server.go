package plotter

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/marksaravi/drone-go/utils"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	SERVER_PORT   int = 8081
	UDP_PORT      int = 6431
	IMU_DATA_SIZE int = 44
	ROTATION_SIZE int = 12
	TIME_SIZE     int = 8
)

func Start() {

	dataChannel := make(chan string)
	startUDPReceiverServer(dataChannel)
	handler := createWebSocketHandler(dataChannel)
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
	time.AfterFunc(time.Second, func() {
		log.Println("Press ENTER to stop server")
	})
	fmt.Scanln()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func createWebSocketHandler(dataChannel chan string) func(w http.ResponseWriter, r *http.Request) {
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

func startUDPReceiverServer(dataChannel chan string) {
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
		log.Printf("UDP Server Listening %s\n", conn.LocalAddr().String())

		const bufferSize int = 16348
		data := make([]byte, bufferSize)

		for {
			nBytes, _, err := conn.ReadFromUDP(data)
			if err != nil {
				panic(err)
			}
			if nBytes > 0 {
				dataPerPacket := int(data[1])
				jsonData := extractPackets(data[2:], dataPerPacket)
				dataChannel <- jsonData

			}
			value += 0.1
		}
	}()
}

func extractPackets(data []byte, dataPerPacket int) string {
	var jsonData string = "["
	var comma string = ""
	for i := 0; i < dataPerPacket; i++ {
		imudata := extractImuRotations(data[IMU_DATA_SIZE*i : IMU_DATA_SIZE*(i+1)])
		jsonData += comma + imudata
		comma = ","
	}
	return jsonData + "]"
}

func extractImuRotations(data []byte) string {
	a := extractRotations(data[0:ROTATION_SIZE])
	g := extractRotations(data[ROTATION_SIZE : 2*ROTATION_SIZE])
	r := extractRotations(data[2*ROTATION_SIZE : 3*ROTATION_SIZE])
	t := utils.UInt64FromBytes(utils.SliceToArray8(data[3*ROTATION_SIZE : 3*ROTATION_SIZE+TIME_SIZE]))
	return fmt.Sprintf("{\"a\":%s,\"g\":%s,\"r\":%s,\"t\":%d}", a, g, r, t)
}

func extractRotations(data []byte) string {
	roll := float64(utils.Float32FromBytes(utils.SliceToArray4(data[0:4])))
	pitch := float64(utils.Float32FromBytes(utils.SliceToArray4(data[4:8])))
	yaw := float64(utils.Float32FromBytes(utils.SliceToArray4(data[8:12])))
	return fmt.Sprintf("{\"roll\":%0.2f,\"pitch\":%0.2f,\"yaw\":%0.2f}", roll, pitch, yaw)
}
