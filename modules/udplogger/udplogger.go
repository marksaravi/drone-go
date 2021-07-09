/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"fmt"
	"net"
	"strings"

	"github.com/MarkSaravi/drone-go/types"
)

const BUFFER_SIZE int = 40

type udpLogger struct {
	conn                 *net.UDPConn
	address              *net.UDPAddr
	enabled              bool
	buffer               []string
	imuDataPerSecond     int
	dataPerPacket        int
	dataPerPacketCounter int
	dataOffset           int
	dataOffsetCounter    int
}

func (l *udpLogger) Enabled() bool {
	return l.enabled
}

func (l *udpLogger) Append(imuRotations types.ImuRotations) {
	if !l.enabled {
		return
	}
	l.dataOffsetCounter++
	if l.dataOffsetCounter == l.dataOffset {
		l.buffer[l.dataPerPacketCounter] = ImuDataToJson(imuRotations)
		l.dataOffsetCounter = 0
		l.dataPerPacketCounter++
	}
}

func (l *udpLogger) Send(imuRotations types.ImuRotations) {
	if !l.Enabled() {
		return
	}
	l.Append(imuRotations)
	l.SendData()
}

func (l *udpLogger) SendData() {
	if l.enabled && l.dataPerPacketCounter == l.dataPerPacket {
		jsonPayload := fmt.Sprintf("{\"d\":[%s],\"dps\":%d}",
			strings.Join(l.buffer[0:l.dataPerPacketCounter], ","),
			l.imuDataPerSecond,
		)
		l.dataPerPacketCounter = 0
		// data len should be less than sysctl net.inet.udp.maxdgram for macOS
		go func() {
			l.conn.WriteToUDP([]byte(jsonPayload), l.address)
		}()
	}
}

func CreateUdpLogger(udpConfig types.UdpLoggerConfig, imuDataPerSecond int) types.UdpLogger {
	if !udpConfig.Enabled {
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
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpConfig.IP, udpConfig.Port))
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return &udpLogger{
			enabled: false,
		}
	}

	var dataPerPacket int = 0
	if udpConfig.PacketsPerSecond > 0 {
		dataPerPacket = udpConfig.DataPerSecond / udpConfig.PacketsPerSecond
		if dataPerPacket > BUFFER_SIZE {
			fmt.Println("Data per packet is more than buffer size BUFFER_SIZE")
			return &udpLogger{
				enabled: false,
			}
		}
	}
	dataOffset := imuDataPerSecond / udpConfig.DataPerSecond
	return &udpLogger{
		conn:                 conn,
		address:              address,
		enabled:              true,
		imuDataPerSecond:     imuDataPerSecond,
		dataPerPacket:        dataPerPacket,
		dataPerPacketCounter: 0,
		dataOffset:           dataOffset,
		dataOffsetCounter:    0,
		buffer:               make([]string, BUFFER_SIZE),
	}
}

func ImuDataToJson(imuRotations types.ImuRotations) string {
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
