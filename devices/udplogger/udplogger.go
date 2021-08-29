/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/MarkSaravi/drone-go/models"
)

type udpLogger struct {
	conn                 *net.UDPConn
	address              *net.UDPAddr
	enabled              bool
	buffer               []string
	imuDataPerSecond     int
	dataPerPacket        int
	dataPerPacketCounter int
	skipOffset           int
	// maxDataPerPacket     int
	bufferCounter int
}

func NewUdpLogger(
	enabled bool,
	ip string,
	port int,
	packetsPerSecond int,
	maxDataPerPacket int,
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

	if packetsPerSecond <= 0 {
		return &udpLogger{
			enabled: false,
		}
	}
	wantedDataPerPacket := imuDataPerSecond / packetsPerSecond
	var skipOffset int = 1
	var actualDataPerPacket = wantedDataPerPacket
	if wantedDataPerPacket > maxDataPerPacket {
		actualDataPerPacket = maxDataPerPacket
		skipOffset = int(math.Ceil(float64(wantedDataPerPacket) / float64(maxDataPerPacket)))
	}
	fmt.Println("DPP: ", imuDataPerSecond, wantedDataPerPacket, actualDataPerPacket, skipOffset)

	return &udpLogger{
		conn:                 conn,
		address:              address,
		enabled:              true,
		imuDataPerSecond:     imuDataPerSecond,
		dataPerPacket:        actualDataPerPacket,
		dataPerPacketCounter: 0,
		bufferCounter:        0,
		skipOffset:           skipOffset,
		buffer:               make([]string, actualDataPerPacket),
	}
}

func (l *udpLogger) appendData(imuRotations models.ImuRotations) {
	if !l.enabled {
		return
	}
	l.dataPerPacketCounter++
	if l.dataPerPacketCounter%l.skipOffset == 0 {
		l.buffer[l.bufferCounter] = imuDataToJson(imuRotations)
		l.bufferCounter++
	}
}

func (l *udpLogger) Send(imuRotations models.ImuRotations) {
	if !l.enabled {
		return
	}
	l.appendData(imuRotations)
	l.sendData()
}

func (l *udpLogger) sendData() {
	if l.enabled && l.dataPerPacketCounter == l.dataPerPacket {
		jsonPayload := fmt.Sprintf("{\"d\":[%s],\"dps\":%d}",
			strings.Join(l.buffer[0:l.bufferCounter], ","),
			l.imuDataPerSecond,
		)
		l.dataPerPacketCounter = 0
		l.bufferCounter = 0
		// data len should be less than sysctl net.inet.udp.maxdgram for macOS
		go func() {
			l.conn.WriteToUDP([]byte(jsonPayload), l.address)
		}()
	}
}

func imuDataToJson(imuRotations models.ImuRotations) string {
	return fmt.Sprintf(`{"a":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"g":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"r":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"t":%d,"dt":%d}`,
		imuRotations.Accelerometer.Roll,
		imuRotations.Accelerometer.Pitch,
		imuRotations.Accelerometer.Yaw,
		imuRotations.Gyroscope.Roll,
		imuRotations.Gyroscope.Pitch,
		imuRotations.Gyroscope.Yaw,
		imuRotations.Rotations.Roll,
		imuRotations.Rotations.Pitch,
		imuRotations.Rotations.Yaw,
		imuRotations.ReadTime,
		imuRotations.ReadInterval,
	)
}
