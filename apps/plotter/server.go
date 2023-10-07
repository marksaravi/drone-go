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

	"github.com/marksaravi/drone-go/utils"
	"nhooyr.io/websocket"
)

const (
	MAX_BUFFER_SIZE     = 8192
	IMU_DATA_SIZE   int = 26
	TIME_SIZE       int = 8
)

type plotter struct {
	udpServerAddress  string
	httpServerAddress string
	httpServer        http.Server
	websocketConn     *websocket.Conn
}

type PlotterSettings struct {
	UDPServerAddress  string
	HTTPServerAddress string
}

func NewPlotter(settings PlotterSettings) *plotter {

	http.Handle("/", http.FileServer(http.Dir("./apps/plotter/static")))
	return &plotter{
		udpServerAddress:  settings.UDPServerAddress,
		httpServerAddress: settings.HTTPServerAddress,
		httpServer: http.Server{
			Addr: settings.HTTPServerAddress,
		},
		websocketConn: nil,
	}
}

func (p *plotter) StartPlotter() {
	fmt.Println("Plotter Started...")
	var waitGroup sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Printf("Shutting Down  HTTP Server...")
		// if err := p.httpServer.Shutdown(context.Background()); err != nil {
		// 	log.Printf("HTTP Server Shutdown Error: %v", err)
		// }
		cancel()
	}()
	p.startUDPServer(ctx, &waitGroup)

	// fmt.Println("HTTP server started...")
	// if err := p.httpServer.ListenAndServe(); err != http.ErrServerClosed {
	// 	log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	// 	cancel()
	// }

	<-ctx.Done()
	waitGroup.Wait()
	fmt.Println("Plotter stopped.")
}

func (p *plotter) startUDPServer(ctx context.Context, wg *sync.WaitGroup) {
	udpAddr, err := net.ResolveUDPAddr("udp", "192.168.1.101:8010")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// // Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	for {
		buf := make([]byte, 8912)
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil && strings.Contains(err.Error(), "closed network connection") {
			fmt.Println(err)
			return
		}

		if err == nil && n > 0 {
			fmt.Println(n)
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// func (p *plotter) startUDPServer(ctx context.Context, wg *sync.WaitGroup) {
// 	wg.Add(1)
// 	fmt.Println("UDP server started...")
// 	go func() {
// 		defer wg.Done()
// 		defer fmt.Println("UDP server stopped.")

// 		address := &net.UDPAddr{
// 			IP:   net.ParseIP("192.168.1.101"),
// 			Port: 8009,
// 		}
// 		udpConn, err := net.ListenUDP("udp", address)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		defer udpConn.Close()
// 		udpConn.SetReadBuffer(8192)
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				err := udpConn.Close()
// 				log.Printf("Closing UDP Connection, error:%v\n", err)
// 				return
// 			default:
// 				udpBuffer := make([]byte, 8192)
// 				udpConn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
// 				n, _, err := udpConn.ReadFromUDP(udpBuffer)
// 				if err != nil && strings.Contains(err.Error(), "closed network connection") {
// 					return
// 				}
// 				if err == nil {
// 					fmt.Println(n)
// 				}
// 			}
// 		}
// 	}()
// }

func (p *plotter) acceptingWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	var err error
	p.websocketConn, err = websocket.Accept(w, r, nil)
	if err != nil {
		p.websocketConn = nil
		log.Println(err)
	} else {
		log.Println("Websocket Connection Accepted")
	}
}

// func Start() {
// 	var wg sync.WaitGroup
// 	ctx, cancel := context.WithCancel(context.Background())

// 	dataChannel := make(chan string)
// 	startUDPReceiverServer(ctx, &wg, dataChannel)
// 	http.Handle("/", http.FileServer(http.Dir("./apps/plotter/static")))
// 	handler := createWebSocketHandler(ctx, &wg, dataChannel)

// 	http.HandleFunc("/conn", handler)
// 	server := http.Server{
// 		Addr: fmt.Sprintf(":%d", SERVER_PORT),
// 	}

// 	go func() {
// 		log.Println("Press ENTER to stop server")
// 		fmt.Scanln()
// 		fmt.Println("Stopping All Servers...")
// 		cancel()
// 	}()
// 	go func() {
// 		log.Println(fmt.Sprintf("Server is listening on port %d\n", SERVER_PORT))
// 		if err := server.ListenAndServe(); err != nil {
// 			log.Printf("HTTP server Shutdown: %v", err)
// 		}
// 	}()
// 	wg.Wait()
// 	fmt.Println("Stopping HTTP Server...")
// 	if err := server.Shutdown(ctx); err != nil {
// 		log.Fatal(err)
// 	}

// }

// func createWebSocketHandler(ctx context.Context, wg *sync.WaitGroup, dataChannel chan string) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Establishing connection")
// 		c, err := websocket.Accept(w, r, nil)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		go func() {
// 			defer wg.Done()

// 			for {
// 				select {
// 				case <-ctx.Done():
// 					c.Close(websocket.StatusInternalError, "closing the websocket connection...")
// 					fmt.Println("Stopping WebSocketHandler...")
// 					return
// 				case json := <-dataChannel:
// 					// fmt.Println(json[0:3])
// 					err = wsjson.Write(ctx, c, json)
// 					if err != nil {
// 						log.Println(err)
// 					}
// 				default:
// 				}
// 			}
// 		}()
// 	}
// }

// func startUDPReceiverServer(ctx context.Context, wg *sync.WaitGroup, dataChannel chan string) {
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		defer fmt.Println("End...")

// 		conn, err := net.ListenUDP("udp", &net.UDPAddr{
// 			Port: UDP_PORT,
// 			IP:   net.ParseIP("0.0.0.0"),
// 		})
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		log.Printf("UDP Server Listening %s\n", conn.LocalAddr().String())

// 		const bufferSize int = 8192
// 		data := make([]byte, bufferSize)

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				close(dataChannel)
// 				conn.Close()
// 				fmt.Println("Stopping UDP Server...")
// 				return
// 			default:
// 				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
// 				nBytes, _, err := conn.ReadFromUDP(data)
// 				if err == nil && nBytes > 0 {
// 					packetSize := int(utils.DeSerializeInt(data[0:2]))
// 					data = data[0:packetSize]
// 					dataPerPacket := int(utils.DeSerializeInt(data[2:4]))
// 					dataLen := int(utils.DeSerializeInt(data[4:6]))
// 					fmt.Println(packetSize, dataPerPacket, dataLen)
// 					jsonData := extractPackets(data[6:], dataLen, dataPerPacket)
// 					// fmt.Println(jsonData[0:32])

// 					dataChannel <- jsonData
// 				}
// 			}
// 		}
// 	}()
// }

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
