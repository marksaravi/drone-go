package utils

import (
	"github.com/marksaravi/drone-go/models"
)

const RADIO_PAYLOAD_SIZE int = 32

func SerializeFlightCommand(flightCommands models.FlightCommands) []byte {
	payload := append([]byte{}, UInt32ToBytes(flightCommands.Id)...)
	payload = append(payload, Float32ToBytes(flightCommands.Roll)...)
	payload = append(payload, Float32ToBytes(flightCommands.Pitch)...)
	payload = append(payload, Float32ToBytes(flightCommands.Yaw)...)
	payload = append(payload, Float32ToBytes(flightCommands.Throttle)...)
	bottons := BoolArrayToByte([8]bool{
		flightCommands.ButtonFrontLeft,
		flightCommands.ButtonFrontRight,
		flightCommands.ButtonTopLeft,
		flightCommands.ButtonTopRight,
		flightCommands.ButtonBottomLeft,
		flightCommands.ButtonBottomRight,
		false,
		false,
	})
	payload = append(payload, ([]byte{bottons, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})...)
	return payload[0:32]
}

func DeserializeFlightCommand(payload []byte) models.FlightCommands {
	buttons := BoolArrayFromByte(payload[20:21][0])
	return models.FlightCommands{
		Id:                UInt32FromBytes(payload[0:4]),
		Roll:              Float32FromBytes(payload[4:8]),
		Pitch:             Float32FromBytes(payload[8:12]),
		Yaw:               Float32FromBytes(payload[12:16]),
		Throttle:          Float32FromBytes(payload[16:20]),
		ButtonFrontLeft:   buttons[0],
		ButtonFrontRight:  buttons[1],
		ButtonTopLeft:     buttons[2],
		ButtonTopRight:    buttons[3],
		ButtonBottomLeft:  buttons[4],
		ButtonBottomRight: buttons[5],
	}
}
