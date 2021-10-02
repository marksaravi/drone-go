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
}

func NewUdpLogger(
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
	}
	logger.buffer.WriteByte(byte(packetsPerSecond))
	logger.buffer.WriteByte(byte(dataPerPacket))
	return &logger
}

func (l *udpLogger) Send(imuRotations models.ImuRotations) {
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

func Append4(buffer *bytes.Buffer, d [4]byte) {
	buffer.Write(d[:])
}

func Append8(buffer *bytes.Buffer, d [8]byte) {
	buffer.Write(d[:])
}

func imuDataToBytes(imuRot models.ImuRotations) []byte {
	buffer := bytes.Buffer{}
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Roll))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Pitch))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Yaw))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Roll))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Pitch))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Yaw))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Roll))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Pitch))
	Append4(&buffer, utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Yaw))
	Append8(&buffer, utils.UInt64ToBytes(uint64(imuRot.ReadTime.UnixNano())))
	return buffer.Bytes()
}

func NewLogger(ctx context.Context, wg *sync.WaitGroup) chan<- models.ImuRotations {
	flightControl := config.ReadFlightControlConfig()
	loggerConfig := config.ReadLoggerConfig()
	loggerConfigs := loggerConfig.UdpLoggerConfigs
	udplogger := NewUdpLogger(
		loggerConfigs.Enabled,
		loggerConfigs.IP,
		loggerConfigs.Port,
		loggerConfigs.PacketsPerSecond,
		flightControl.Configs.ImuDataPerSecond,
	)
	loggerChan := make(chan models.ImuRotations)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer log.Println("LOGGER CLOSED")
		defer close(loggerChan)

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-loggerChan:
				udplogger.Send(data)
			}
		}
	}()
	return loggerChan
}
