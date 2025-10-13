package plotter

import (
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

func SerializeHeader() []byte {
	packet := make([]byte, 0, PLOTER_PACKET_HEADER_LEN)
	packet = append(packet, utils.SerializeInt(PLOTTER_PACKET_LEN)...)
	packet = append(packet, utils.SerializeInt(PLOTTER_DATA_PER_PACKET)...)
	packet = append(packet, utils.SerializeInt(PLOTTER_DATA_LEN)...)
	return packet
}

func DeSerializeHeader(packet []byte) (packetSize, dataPerPacket, dataLen int) {
	packetSize = int(utils.DeSerializeInt(packet[0:PLOTTER_INT_DATA_LEN]))
	// fmt.Println("packetSize:", packetSize)
	dataPerPacket = int(utils.DeSerializeInt(packet[PLOTTER_INT_DATA_LEN : 2*PLOTTER_INT_DATA_LEN]))
	// fmt.Println("dataPerPacket:", dataPerPacket)
	dataLen = int(utils.DeSerializeInt(packet[PLOTTER_INT_DATA_LEN*2 : PLOTTER_INT_DATA_LEN*3]))
	return
}

func SerializeDroneData(dur time.Duration, rotations, accelerometer, gyroscope imu.Rotations, throttle byte) []byte {
	packet := make([]byte, 0, PLOTTER_DATA_LEN)

	packet = append(packet, utils.SerializeDuration(dur)...)
	packet = append(packet, utils.SerializeFloat64(rotations.Roll)...)
	packet = append(packet, utils.SerializeFloat64(rotations.Pitch)...)
	packet = append(packet, utils.SerializeFloat64(rotations.Yaw)...)
	packet = append(packet, utils.SerializeFloat64(accelerometer.Roll)...)
	packet = append(packet, utils.SerializeFloat64(accelerometer.Pitch)...)
	packet = append(packet, utils.SerializeFloat64(accelerometer.Yaw)...)
	packet = append(packet, utils.SerializeFloat64(gyroscope.Roll)...)
	packet = append(packet, utils.SerializeFloat64(gyroscope.Pitch)...)
	packet = append(packet, utils.SerializeFloat64(gyroscope.Yaw)...)
	packet = append(packet, throttle)
	return packet
}

func DeSerializeDroneData(dataPacket []byte) (dur time.Duration, rotations, accelerometer, gyroscope imu.Rotations, throttle byte) {
	dur = utils.DeSerializeDuration(dataPacket[0:PLOTTER_DUR_DATA_LEN])
	floats := make([]float64, 9)
	for i := 0; i < 9; i++ {
		floats[i] = utils.DeSerializeFloat64(dataPacket[PLOTTER_DUR_DATA_LEN+PLOTTER_FLOAT_DATA_LEN*i : PLOTTER_DUR_DATA_LEN+PLOTTER_FLOAT_DATA_LEN*(i+1)])
	}
	rotations = imu.Rotations{
		Roll:  floats[0],
		Pitch: floats[1],
		Yaw:   floats[2],
	}
	accelerometer = imu.Rotations{
		Roll:  floats[3],
		Pitch: floats[4],
		Yaw:   floats[5],
	}
	gyroscope = imu.Rotations{
		Roll:  floats[6],
		Pitch: floats[7],
		Yaw:   floats[8],
	}
	throttle = dataPacket[PLOTTER_DUR_DATA_LEN+PLOTTER_FLOAT_DATA_LEN*9]
	return
}
