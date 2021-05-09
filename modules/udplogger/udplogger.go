package udplogger

import (
	"fmt"
	"net"

	"github.com/MarkSaravi/drone-go/types"
)

type udpLogger struct {
	con       *net.PacketConn
	address   *net.UDPAddr
	enabled   bool
	buffer    []string
	bufferLen int
}

func (l *udpLogger) Send(json string) {
	if !l.enabled {
		return
	}
	l.buffer = append(l.buffer, json)
	if len(l.buffer) == l.bufferLen {
		data := "["
		comma := ""
		for _, s := range l.buffer {
			data = data + comma + s
			comma = ","
		}
		data = data + "]"
		l.buffer = nil
		go func() {
			(*l.con).WriteTo([]byte(data), l.address)
		}()
	}
}

func CreateUdpLogger(
	udpConfig types.UdpLoggerConfig,
	dataPerSecond int) udpLogger {
	if !udpConfig.Enabled {
		fmt.Println("UDP is not enabled")
		return udpLogger{
			enabled: false,
		}
	}
	con, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return udpLogger{
			enabled: false,
		}
	}
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpConfig.IP, udpConfig.Port))
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return udpLogger{
			enabled: false,
		}
	}
	var bufferLen int = dataPerSecond / 50
	if bufferLen == 0 {
		bufferLen = 1
	}
	return udpLogger{
		con:       &con,
		address:   address,
		enabled:   true,
		bufferLen: bufferLen,
	}
}
