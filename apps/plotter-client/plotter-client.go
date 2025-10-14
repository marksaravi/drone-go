package plotterclient

import (
	"fmt"
	"net"
	"time"

	"github.com/marksaravi/drone-go/apps/plotter"
	"github.com/marksaravi/drone-go/devices/imu"
)

const PLOTTER_DATA_PER_PACKET = 128
const PLOTTER_PACKET_LEN = plotter.PLOTER_PACKET_HEADER_LEN + PLOTTER_DATA_PER_PACKET*plotter.PLOTTER_DATA_LEN

type plotterClient struct {
	active        bool
	address       string
	dataPacket    []byte
	sendBuffer    []byte
	dataCounter   int
	dataPerPacket int
	startTime     time.Time
	udpConn       *net.UDPConn
}

type Settings struct {
	Active  bool
	Address string
}

func NewPlotter(settings Settings) *plotterClient {
	p := plotterClient{
		active:        settings.Active,
		address:       settings.Address,
		dataPacket:    make([]byte, 0, PLOTTER_PACKET_LEN),
		sendBuffer:    make([]byte, PLOTTER_PACKET_LEN),
		dataCounter:   0,
		dataPerPacket: PLOTTER_DATA_PER_PACKET,
	}
	if p.active {
		p.initUdp()
	}
	return &p
}

func (p *plotterClient) initUdp() {
	plotterUdpServer, err := net.ResolveUDPAddr("udp", p.address)
	if err != nil {
		p.active = false
		fmt.Println("unable to initialise plotter server. Plotter deactivated.")
	}
	p.udpConn, err = net.DialUDP("udp", nil, plotterUdpServer)
	if err != nil || p.udpConn == nil {
		p.active = false
		fmt.Println("unable to initialise plotter connection. Plotter deactivated.")
	}
}

func (p *plotterClient) SendPlotterData(rotations, accRotations, gyroRotations imu.Rotations) bool {
	if !p.active {
		return false
	}
	if p.dataCounter == 0 {
		p.dataPacket = make([]byte, 0, PLOTTER_PACKET_LEN)
		p.dataPacket = append(p.dataPacket, plotter.SerializeHeader(PLOTTER_PACKET_LEN)...)
	}
	p.SerializeRotations(rotations, accRotations, gyroRotations)
	if p.dataCounter < p.dataPerPacket {
		return false
	}
	if p.udpConn != nil {
		copy(p.sendBuffer, p.dataPacket)
		go func() {
			p.udpConn.Write(p.sendBuffer)
		}()
	}
	p.dataCounter = 0
	return true
}

func (p *plotterClient) SetStartTime(startTime time.Time) {
	p.startTime = startTime
}

func (p *plotterClient) SerializeRotations(rotations, accRotations, gyroRotations imu.Rotations) {
	p.dataPacket = append(
		p.dataPacket,
		plotter.SerializeDroneData(
			time.Since(p.startTime),
			rotations,
			accRotations,
			gyroRotations,
			0,
		)...)
	p.dataCounter++
}
