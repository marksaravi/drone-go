package utils

import (
	"bytes"

	"github.com/marksaravi/drone-go/models"
)

const RADIO_PAYLOAD_SIZE int = 32

func SerializeFlightCommand(flightCommands models.FlightCommands) [32]byte {
	var buf bytes.Buffer
	typeBytes := [1]byte{flightCommands.PayloadType}
	idBytes := UInt32ToBytes(flightCommands.Id)
	timeBytes := Int64ToBytes(flightCommands.Time)
	rollBytes := Float32ToBytes(flightCommands.Roll)
	pitchBytes := Float32ToBytes(flightCommands.Pitch)
	yawBytes := Float32ToBytes(flightCommands.Yaw)
	throttleBytes := Float32ToBytes(flightCommands.Throttle)
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
	buf.Write(typeBytes[:])
	buf.Write(idBytes[:])
	buf.Write(timeBytes[:])
	buf.Write(rollBytes[:])
	buf.Write(pitchBytes[:])
	buf.Write(yawBytes[:])
	buf.Write(throttleBytes[:])
	buf.WriteByte(bottons)

	var payload [32]byte
	copy(payload[:], buf.Bytes())
	return payload
}

func DeserializeFlightCommand(payload [32]byte) models.FlightCommands {
	buttons := BoolArrayFromByte(payload[29])
	return models.FlightCommands{
		PayloadType:       payload[0],
		Id:                UInt32FromBytes(SliceToArray4(payload[1:5])),
		Time:              Int64FromBytes(SliceToArray8(payload[5:13])),
		Roll:              Float32FromBytes(SliceToArray4(payload[13:17])),
		Pitch:             Float32FromBytes(SliceToArray4(payload[17:21])),
		Yaw:               Float32FromBytes(SliceToArray4(payload[21:25])),
		Throttle:          Float32FromBytes(SliceToArray4(payload[25:29])),
		ButtonFrontLeft:   buttons[0],
		ButtonFrontRight:  buttons[1],
		ButtonTopLeft:     buttons[2],
		ButtonTopRight:    buttons[3],
		ButtonBottomLeft:  buttons[4],
		ButtonBottomRight: buttons[5],
	}
}
