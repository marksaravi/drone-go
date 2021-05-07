package main

import (
	"fmt"
	"net"
)

const LOGGER_BUFFER_LEN = 1

type UdpLogger struct {
	con           *net.PacketConn
	address       *net.UDPAddr
	enabled       bool
	buffer        [LOGGER_BUFFER_LEN]string
	bufferCounter int
}

func (l *UdpLogger) send(json string) {
	if !l.enabled {
		return
	}
	l.buffer[l.bufferCounter] = json
	l.bufferCounter += 1
	if l.bufferCounter == LOGGER_BUFFER_LEN {
		data := ""
		for _, s := range l.buffer {
			data = data + "\n" + s
		}
		l.bufferCounter = 0
		go func() {
			(*l.con).WriteTo([]byte(data), l.address)
		}()
	}
}

func createUdpConnection(appConfig ApplicationConfig) UdpLogger {
	if !appConfig.UDP.Enabled {
		fmt.Println("UDP is not enabled")
		return UdpLogger{
			enabled: false,
		}
	}
	con, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return UdpLogger{
			enabled: false,
		}
	}
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", appConfig.UDP.IP, appConfig.UDP.Port))
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return UdpLogger{
			enabled: false,
		}
	}
	return UdpLogger{
		con:           &con,
		address:       address,
		enabled:       true,
		bufferCounter: 0,
	}
}
