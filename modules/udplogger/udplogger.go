package udplogger

import (
	"fmt"
	"net"

	"github.com/MarkSaravi/drone-go/types"
)

type udpLogger struct {
	conn          *net.UDPConn
	address       *net.UDPAddr
	enabled       bool
	buffer        []string
	bufferLen     int
	dataPerSecond int
	sendFrequency int
	maxDataSize   int
	chunkId       byte
}

func (l *udpLogger) Send(json string) {
	if !l.enabled {
		return
	}
	l.buffer = append(l.buffer, json)
	if len(l.buffer) == l.bufferLen {
		jsonArray := ""
		comma := ""
		for _, s := range l.buffer {
			jsonArray = jsonArray + comma + s
			comma = ","
		}
		data := fmt.Sprintf("{\"data\": [%s], \"dataPerSecond\": %d, \"packetsPerSecond\": %d}",
			jsonArray,
			l.dataPerSecond,
			l.sendFrequency,
		)
		l.buffer = nil
		go func() {
			bytes := []byte(data)
			total := len(bytes)
			var numOfChunks byte = byte(total / l.maxDataSize)
			if total%l.maxDataSize > 0 {
				numOfChunks += 1
			}
			var chunkCounter byte = 0
			var chunk []byte
			for from := 0; from < total; from += l.maxDataSize {
				if from+l.maxDataSize <= total {
					chunk = bytes[from : from+l.maxDataSize]
				} else {
					chunk = bytes[from:]
				}
				chunkInfo := []byte{l.chunkId, chunkCounter, numOfChunks}
				udpdata := append(chunkInfo, chunk...)
				l.conn.WriteToUDP(udpdata, l.address)
				chunkCounter += 1
			}
			l.chunkId += 1
			if l.chunkId > 100 {
				l.chunkId = 0
			}
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
	conn, err := net.ListenUDP("udp", nil)
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
	const sendFrequency = 50
	var bufferLen int = dataPerSecond / sendFrequency
	if bufferLen == 0 {
		bufferLen = 1
	}
	return udpLogger{
		conn:          conn,
		address:       address,
		enabled:       true,
		bufferLen:     bufferLen,
		dataPerSecond: dataPerSecond,
		sendFrequency: sendFrequency,
		maxDataSize:   9000,
		chunkId:       0,
	}
}
