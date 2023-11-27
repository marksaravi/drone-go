package drone

import (
	"fmt"
	"net"
	"time"

	"github.com/marksaravi/drone-go/apps/plotter"
)

// func (d *droneApp) ReadIMU() bool {
// 	if time.Since(d.lastImuData) < time.Second/time.Duration(d.imuDataPerSecond) {
// 		return false
// 	}
// 	d.lastImuData = time.Now()
// 	rotations, accRotations, gyroRotations, err := d.imu.Read()
// 	if err != nil {
// 		return false
// 	}
// 	d.rotations = rotations
// 	d.accRotations = accRotations
// 	d.gyroRotations = gyroRotations
// 	return true
// }

// func (d *droneApp) ReceiveCommand() ([]byte, bool) {
// 	if time.Since(d.lastCommand) < time.Second/time.Duration(2*d.commandsPerSecond) {
// 		return nil, false
// 	}
// 	d.lastCommand = time.Now()
// 	return d.receiver.Receive()
// }

func (d *droneApp) InitUdp() {
	if !d.plotterActive {
		return
	}
	plotterUdpServer, err := net.ResolveUDPAddr("udp", d.plotterAddress)
	if err != nil {
		d.plotterActive = false
		fmt.Println("unable to initialise plotter server. Plotter deactivated.")
	}
	d.plotterUdpConn, err = net.DialUDP("udp", nil, plotterUdpServer)
	if err != nil || d.plotterUdpConn == nil {
		d.plotterActive = false
		fmt.Println("unable to initialise plotter connection. Plotter deactivated.")
	}
}
func (d *droneApp) SendPlotterData() bool {
	if !d.plotterActive {
		return false
	}
	if d.plotterDataCounter == 0 {
		d.plotterDataPacket = make([]byte, 0, plotter.PLOTTER_PACKET_LEN)
		d.plotterDataPacket = append(d.plotterDataPacket, plotter.SerializeHeader()...)
	}
	d.SerializeRotations()
	if d.plotterDataCounter < d.ploterDataPerPacket {
		return false
	}
	if d.plotterUdpConn != nil {
		copy(d.plotterSendBuffer, d.plotterDataPacket)
		go func() {
			d.plotterUdpConn.Write(d.plotterSendBuffer)
		}()
	}
	d.plotterDataCounter = 0
	return true
}

func (d *droneApp) SerializeRotations() {
	d.plotterDataPacket = append(
		d.plotterDataPacket,
		plotter.SerializeDroneData(
			time.Since(d.startTime),
			d.rotations,
			d.accRotations,
			d.gyroRotations,
			0,
		)...)
	d.plotterDataCounter++
}
