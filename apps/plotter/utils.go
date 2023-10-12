package plotter

import (
	"time"

	"github.com/marksaravi/drone-go/devices/imu"
	"github.com/marksaravi/drone-go/utils"
)

func SerializeHeader() []byte {
	packet:=make([]byte, 0, PLOTER_PACKET_HEADER_SIZE)
	packet=append(packet, utils.SerializeInt(PLOTTER_PACKET_SIZE)...)
	packet=append(packet, utils.SerializeInt(PLOTTER_DATA_PER_PACKET)...)
	packet=append(packet, utils.SerializeInt(PLOTTER_DATA_LEN)...)
	return packet
}

func DeSerializeHeader(packet []byte) (packetSize, dataPerPacket, dataLen int) {
	packetSize = int(utils.DeSerializeInt(packet[0:2]))
	dataPerPacket = int(utils.DeSerializeInt(packet[2:4]))
	dataLen = int(utils.DeSerializeInt(packet[4:6]))
	return
}

func SerializeDroneData(dur time.Duration, rotations, accelerometer, gyroscope imu.Rotations, throttle byte) []byte {
	packet:=make([]byte, 0, PLOTTER_DATA_LEN)

	packet=append(packet, utils.SerializeDuration(dur)...)
	packet=append(packet, utils.SerializeFloat64(rotations.Roll)...)
	packet=append(packet, utils.SerializeFloat64(rotations.Pitch)...)
	packet=append(packet, utils.SerializeFloat64(rotations.Yaw)...)
	packet=append(packet, utils.SerializeFloat64(accelerometer.Roll)...)
	packet=append(packet, utils.SerializeFloat64(accelerometer.Pitch)...)
	packet=append(packet, utils.SerializeFloat64(accelerometer.Yaw)...)
	packet=append(packet, utils.SerializeFloat64(gyroscope.Roll)...)
	packet=append(packet, utils.SerializeFloat64(gyroscope.Pitch)...)
	packet=append(packet, utils.SerializeFloat64(gyroscope.Yaw)...)
	packet=append(packet, throttle)
	return packet
}

func DeSerializeDroneData(dataPacket []byte) (dur time.Duration, rotations, accelerometer, gyroscope imu.Rotations, throttle byte) {
	dur = utils.DeSerializeDuration(dataPacket[0:4])
	floats:=make([]float64, 9)
	for i:=0; i<9; i++ {
		floats[i] = utils.DeSerializeFloat64(dataPacket[4+i*2:4+(i+1)*2])
	}
	rotations = imu.Rotations {
		Roll: utils.DeSerializeFloat64(dataPacket[4:6]),
		Pitch: utils.DeSerializeFloat64(dataPacket[6:8]),
		Yaw: utils.DeSerializeFloat64(dataPacket[8:10]),
	} 
	accelerometer = imu.Rotations {
		Roll: utils.DeSerializeFloat64(dataPacket[10:12]),
		Pitch: utils.DeSerializeFloat64(dataPacket[12:14]),
		Yaw: utils.DeSerializeFloat64(dataPacket[14:16]),
	} 
	gyroscope = imu.Rotations {
		Roll: utils.DeSerializeFloat64(dataPacket[16:18]),
		Pitch: utils.DeSerializeFloat64(dataPacket[18:20]),
		Yaw: utils.DeSerializeFloat64(dataPacket[20:22]),
	}
	throttle=dataPacket[22]
	return
}