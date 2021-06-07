/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/MarkSaravi/drone-go/types"
)

type udpLogger struct {
	conn             *net.UDPConn
	address          *net.UDPAddr
	enabled          bool
	buffer           []string
	dataPerPacket    int
	dataPerSecond    int
	packetsPerSecond int
	printIntervalMs  int
	lastPrint        time.Time
}

func (l *udpLogger) Send(json string) {
	if !l.enabled || l.packetsPerSecond == 0 {
		return
	}
	l.buffer = append(l.buffer, json)
	if time.Since(l.lastPrint) >= time.Duration(time.Millisecond*time.Duration(l.printIntervalMs)) {
		l.lastPrint = time.Now()
	}
	if len(l.buffer) == l.dataPerPacket {
		jsonPayload := fmt.Sprintf("{\"d\":[%s],\"dps\":%d,\"pps\":%d,\"dpp\":%d}",
			strings.Join(l.buffer, ","),
			l.dataPerSecond,
			l.packetsPerSecond,
			l.dataPerPacket,
		)
		l.buffer = nil
		// data len should be less than sysctl net.inet.udp.maxdgram for macOS
		go func() {
			l.conn.WriteToUDP([]byte(jsonPayload), l.address)
		}()
	}
}

func CreateUdpLogger(
	udpConfig types.UdpLoggerConfig,
	dataPerSecond int,
) udpLogger {
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

	var dataPerPacket int = 0
	if udpConfig.PacketsPerSecond > 0 {
		dataPerPacket = dataPerSecond / udpConfig.PacketsPerSecond
	}
	return udpLogger{
		conn:             conn,
		address:          address,
		enabled:          true,
		dataPerPacket:    dataPerPacket,
		dataPerSecond:    dataPerSecond,
		packetsPerSecond: udpConfig.PacketsPerSecond,
		printIntervalMs:  udpConfig.PrintIntervalMs,
		lastPrint:        time.Now(),
	}
}

// https://github.com/google/gopacket
// func (s *scanner) sendUDPHW(dstIP string, idx int, ipPtr *layers.IPv4, udpPtr *layers.UDP, buf gopacket.SerializeBuffer) {
// 	srcPort = s.srcPorts[idx]
// 	pldPtr = &s.payloads[idx]
// 	ipPtr.DstIP = net.ParseIP(dstIP)
// 	udpPtr.SrcPort = layers.UDPPort(srcPort)

// 	if err := gopacket.SerializeLayers(buf, s.opts, &s.eth, ipPtr, udpPtr, pldPtr); err != nil {...}
// 	s.handle.WritePacketData(buf.Bytes())
// }
