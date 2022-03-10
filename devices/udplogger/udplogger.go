/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"bytes"
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

func (l *udpLogger) sendUDP(imuRotations models.ImuRotations) {
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

func NewUdpLogger(loggerConfigs config.UdpLoggerConfigs, dataPerSecond int) *udpLogger {
	var conn *net.UDPConn = nil
	var err error = nil
	var address *net.UDPAddr = nil
	var enabled bool = false

	if loggerConfigs.Enabled {
		enabled = loggerConfigs.Enabled
		conn, err = net.ListenUDP("udp", nil)
		if err != nil {
			fmt.Println("UDP initialization error: ", err)
			enabled = false
		}
		address, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", loggerConfigs.IP, loggerConfigs.Port))
		if err != nil {
			fmt.Println("UDP initialization error: ", err)
			enabled = false
		}
	}

	dataPerPacket := dataPerSecond / loggerConfigs.PacketsPerSecond

	logger := udpLogger{
		enabled:          enabled,
		conn:             conn,
		address:          address,
		dataPerPacket:    dataPerPacket,
		packetsPerSecond: loggerConfigs.PacketsPerSecond,
		bufferCounter:    0,
		buffer:           bytes.Buffer{},
		dataChannel:      make(chan models.ImuRotations),
	}
	logger.buffer.WriteByte(byte(loggerConfigs.PacketsPerSecond))
	logger.buffer.WriteByte(byte(dataPerPacket))
	return &logger

}

func (l *udpLogger) Close() {
	close(l.dataChannel)
}

func (l *udpLogger) Send(data models.ImuRotations) {
	if l.enabled {
		l.dataChannel <- data
	}
}

func (l *udpLogger) Start(wg *sync.WaitGroup) {
	if !l.enabled {
		log.Println("Logger is not enabled.")
		return
	}
	wg.Add(1)

	log.Println("Starting the Logger...")
	go func() {
		defer wg.Done()
		defer log.Println("Loger is stopped.")
		var loggerChanOpen bool = true
		var imuData models.ImuRotations
		for loggerChanOpen {
			imuData, loggerChanOpen = <-l.dataChannel
			if loggerChanOpen {
				l.sendUDP(imuData)
			}
		}
	}()
}
