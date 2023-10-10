package drone

import (
	"fmt"
	"net"
	"time"

	"github.com/marksaravi/drone-go/constants"
	"github.com/marksaravi/drone-go/utils"
)

func (d *droneApp) ReadIMU() bool {
	if time.Since(d.lastImuData) < time.Second/time.Duration(d.imuDataPerSecond) {
		return false
	}
	d.lastImuData = time.Now()
	rotations, accRotations, gyroRotations, err := d.imu.Read()
	if err != nil {
		return false
	}
	d.rotations = rotations
	d.accRotations = accRotations
	d.gyroRotations = gyroRotations
	return true
}

func (d *droneApp) ReceiveCommand() ([]byte, bool) {
	if time.Since(d.lastCommand) < time.Second/time.Duration(2*d.commandsPerSecond) {
		return nil, false
	}
	d.lastCommand = time.Now()
	return d.receiver.Receive()
}

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
		d.plotterDataPacket = make([]byte, 0, constants.PLOTTER_PACKET_SIZE)
		d.serialiseInt(constants.PLOTTER_PACKET_SIZE)
		d.serialiseInt(constants.PLOTTER_DATA_PER_PACKET)
		d.serialiseInt(constants.PLOTTER_DATA_LEN)
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
	d.serialiseTimeStamp(time.Since(d.startTime))
	d.serialiseRotattion(d.rotations.Roll)
	d.serialiseRotattion(d.rotations.Pitch)
	d.serialiseRotattion(d.rotations.Yaw)
	d.serialiseRotattion(d.accRotations.Roll)
	d.serialiseRotattion(d.accRotations.Pitch)
	d.serialiseRotattion(d.accRotations.Yaw)
	d.serialiseRotattion(d.gyroRotations.Roll)
	d.serialiseRotattion(d.gyroRotations.Pitch)
	d.serialiseRotattion(d.gyroRotations.Yaw)
	d.plotterDataCounter++
}

func (d *droneApp) serialiseTimeStamp(dur time.Duration) {
	d.plotterDataPacket = append(d.plotterDataPacket, utils.SerializeDuration(dur)...)
}

func (d *droneApp) serialiseInt(n int) {
	d.plotterDataPacket = append(d.plotterDataPacket, utils.SerializeInt(int16(n))...)
}

func (d *droneApp) serialiseRotattion(r float64) {
	d.plotterDataPacket = append(d.plotterDataPacket, utils.SerializeFloat64(r)...)
}
