/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"bytes"
	"fmt"
	"net"

	"github.com/MarkSaravi/drone-go/models"
	"github.com/MarkSaravi/drone-go/utils"
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

func imuDataToBytes(imuRot models.ImuRotations) []byte {
	buffer := bytes.Buffer{}
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Roll))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Pitch))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Accelerometer.Yaw))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Roll))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Pitch))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Gyroscope.Yaw))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Roll))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Pitch))
	buffer.Write(utils.Float64ToRoundedFloat32Bytes(imuRot.Rotations.Yaw))
	buffer.Write(utils.UInt64ToBytes(uint64(imuRot.ReadTime.UnixNano())))
	return buffer.Bytes()
}
