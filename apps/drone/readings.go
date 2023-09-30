package drone

import (
	"encoding/binary"
	"fmt"
	"time"
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

func (d *droneApp) SendPlotterData() {
	d.plotterDataCounter++
	if d.plotterDataCounter < PLOTTER_SAMPLE_RATE {
		return
	}
	d.plotterDataCounter = 0
	buffer := make([]byte, len(d.plotterBuffer))
	copy(buffer, d.plotterBuffer)
	d.plotterBuffer = make([]byte, 0, PLOTTER_BUFFER_SIZE)
	go func(b []byte) {
		fmt.Println("sending plotter data")
	}(buffer)
}
func (d *droneApp) BufferPlotterData() {
	if !d.plotterActive {
		return
	}
	d.ploterSampleIntervalCounter++
	if d.ploterSampleIntervalCounter < d.plotterSampleInterval {
		return
	}
	d.ploterSampleIntervalCounter = 0
	d.plotterBuffer = append(d.plotterBuffer, d.SerializeRotations()...)
}

func (d *droneApp) SerializeRotations() []byte {
	data := make([]byte, 0, 64)
	t := make([]byte, 8)
	binary.LittleEndian.PutUint64(t, uint64(time.Now().UnixMicro()))
	data = append(data, t...)
	data = append(data, rotationToInt(d.rotations.Roll)...)
	data = append(data, rotationToInt(d.rotations.Pitch)...)
	data = append(data, rotationToInt(d.rotations.Yaw)...)
	data = append(data, rotationToInt(d.accRotations.Roll)...)
	data = append(data, rotationToInt(d.accRotations.Pitch)...)
	data = append(data, rotationToInt(d.accRotations.Yaw)...)
	data = append(data, rotationToInt(d.gyroRotations.Roll)...)
	data = append(data, rotationToInt(d.gyroRotations.Pitch)...)
	data = append(data, rotationToInt(d.gyroRotations.Yaw)...)

	return data
}

func rotationToInt(r float64) []byte {
	n := uint16(r*10) + 16000
	d := []byte{0, 0}
	binary.LittleEndian.PutUint16(d, n)
	return d
}
