package main

import (
	"fmt"
	"time"

	"github.com/marksaravi/drone-go/hardware"
	"github.com/marksaravi/drone-go/hardware/mcp3008"
)

func main() {
	analogToDigitalSPIConn := hardware.NewSPIConnection(
		1,
		0,
	)
	xAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		0,
		512,
		180,
	)
	yAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		1,
		512,
		180,
	)
	zAxisAnalogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		2,
		512,
		180,
	)
	throttleAlogToDigitalConvertor := mcp3008.NewMCP3008(
		analogToDigitalSPIConn,
		3,
		512,
		200,
	)

	for {
		time.Sleep(time.Second / 2)
		fmt.Println(xAxisAnalogToDigitalConvertor.Read())
		fmt.Println(yAxisAnalogToDigitalConvertor.Read())
		fmt.Println(zAxisAnalogToDigitalConvertor.Read())
		fmt.Println(throttleAlogToDigitalConvertor.Read())
	}
}
