/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/marksaravi/drone-go/config"
	"github.com/marksaravi/drone-go/models"
	"github.com/marksaravi/drone-go/utils"
)

type udpLogger struct {
	conn             *net.UDPConn
	address          *net.UDPAddr
	enabled          bool
	buffer           bytes.Buffer
	packetsPerSecond int
	dataPerPacket    int
	bufferCounter    int
	dataChannel      chan models.ImuRotations
}

func newLogger(
	enabled bool,
	ip string,
	port int,
	packetsPerSecond int,
	imuDataPerSecond int,
) *udpLogger {
	if !enabled {
		return &udpLogger{
			enabled: false,
		}
	}
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return &udpLogger{
			enabled: false,
		}
	}
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return &udpLogger{
			enabled: false,
		}
	}

	dataPerPacket := imuDataPerSecond / packetsPerSecond

	logger := udpLogger{
		conn:             conn,
		address:          address,
		enabled:          true,
		dataPerPacket:    dataPerPacket,
		packetsPerSecond: packetsPerSecond,
		bufferCounter:    0,
		buffer:           bytes.Buffer{},
		dataChannel:      make(chan models.ImuRotations),
	}
	logger.buffer.WriteByte(byte(packetsPerSecond))
	logger.buffer.WriteByte(byte(dataPerPacket))
	return &logger
}

func (l *udpLogger) send(imuRotations models.ImuRotations) {
	if !l.enabled {
		return
	}
	data := imuDataToBytes(imuRotations)
	l.buffer.Write(data)
	l.bufferCounter++
	if l.bufferCounter == l.dataPerPacket {
		payload := l.buffer.Bytes()
		l.conn.WriteToUDP(payload, l.address)
		l.buffer = bytes.Buffer{}
		l.buffer.WriteByte(payload[0])
		l.buffer.WriteByte(payload[1])
		l.bufferCounter = 0
	}

}

func append4Bytes(buffer *bytes.Buffer, d [4]byte) {
	buffer.Write(d[:])
}

func append8Bytes(buffer *bytes.Buffer, d [8]byte) {
	buffer.Write(d[:])
}

func imuDataToBytes(imuRot models.ImuRotations) []byte {
	buffer := bytes.Buffer{}
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Roll))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Pitch))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Yaw))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Roll))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Pitch))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Yaw))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Roll))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Pitch))
	append4Bytes(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Yaw))
	append8Bytes(&buffer, utils.UInt64ToBytes(uint64(imuRot.ReadTime.UnixNano())))
	return buffer.Bytes()
}

func NewUdpLogger() *udpLogger {
	flightControl := config.ReadFlightControlConfig()
	loggerConfig := config.ReadLoggerConfig()
	loggerConfigs := loggerConfig.UdpLoggerConfigs
	return newLogger(
		loggerConfigs.Enabled,
		loggerConfigs.IP,
		loggerConfigs.Port,
		loggerConfigs.PacketsPerSecond,
		flightControl.Configs.ImuDataPerSecond,
	)
}

func (l *udpLogger) Close() {
	close(l.dataChannel)
	l.dataChannel = nil
}

func (l *udpLogger) Send(data models.ImuRotations) {
	if l.dataChannel != nil {
		l.dataChannel <- data
	}
}

func (l *udpLogger) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	log.Println("Starting the Logger...")
	go func() {
		defer wg.Done()
		defer log.Println("Loger is stopped.")
		for l.dataChannel != nil {
			select {
			case data, ok := <-l.dataChannel:
				if ok {
					l.send(data)
				}
			default:
			}
		}
	}()
}
