/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/MarkSaravi/drone-go/models"
)

type udpLogger struct {
	conn          *net.UDPConn
	address       *net.UDPAddr
	enabled       bool
	buffer        bytes.Buffer
	dataPerPacket int
	bufferCounter int
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

	return &udpLogger{
		conn:          conn,
		address:       address,
		enabled:       true,
		dataPerPacket: dataPerPacket,
		bufferCounter: 0,
		buffer:        bytes.Buffer{},
	}
}

func (l *udpLogger) Send(imuRotations models.ImuRotations) {
	if !l.enabled {
		return
	}
	data := imuDataToBytes(imuRotations)
	l.buffer.Write(data)
	l.bufferCounter++
	if l.bufferCounter == l.dataPerPacket {
		fmt.Println(len(l.buffer.Bytes()), l.dataPerPacket, len(data))
		l.conn.WriteToUDP(l.buffer.Bytes(), l.address)
		l.buffer = bytes.Buffer{}
		l.bufferCounter = 0
	}

}

func float64ToTransferBytes(x float64) []byte {
	var i int16 = int16(math.Round(x * 100))
	var ui uint16 = uint16(i + 32767)
	return uint16ToBytes(ui)
}

func imuDataToBytes(imuRot models.ImuRotations) []byte {
	buffer := bytes.Buffer{}
	buffer.Write(float64ToTransferBytes(imuRot.Accelerometer.Roll))
	buffer.Write(float64ToTransferBytes(imuRot.Accelerometer.Pitch))
	buffer.Write(float64ToTransferBytes(imuRot.Accelerometer.Yaw))
	buffer.Write(float64ToTransferBytes(imuRot.Gyroscope.Roll))
	buffer.Write(float64ToTransferBytes(imuRot.Gyroscope.Pitch))
	buffer.Write(float64ToTransferBytes(imuRot.Gyroscope.Yaw))
	buffer.Write(float64ToTransferBytes(imuRot.Rotations.Roll))
	buffer.Write(float64ToTransferBytes(imuRot.Rotations.Pitch))
	buffer.Write(float64ToTransferBytes(imuRot.Rotations.Yaw))
	buffer.Write(uint64ToBytes(uint64(imuRot.ReadTime)))
	return buffer.Bytes()
}

func uint16ToBytes(i uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, i)
	return buf
}

func uint64ToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, i)
	return buf
}
