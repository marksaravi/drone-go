package drone

func (d *droneApp) InitUdp() {
	// if !d.plotterActive {
	// 	return
	// }
	// plotterUdpServer, err := net.ResolveUDPAddr("udp", d.plotterAddress)
	// if err != nil {
	// 	d.plotterActive = false
	// 	fmt.Println("unable to initialise plotter server. Plotter deactivated.")
	// }
	// d.plotterUdpConn, err = net.DialUDP("udp", nil, plotterUdpServer)
	// if err != nil || d.plotterUdpConn == nil {
	// 	d.plotterActive = false
	// 	fmt.Println("unable to initialise plotter connection. Plotter deactivated.")
	// }
}

func (d *droneApp) SendPlotterData() bool {
	// if !d.plotterActive {
	// 	return false
	// }
	// if d.plotterDataCounter == 0 {
	// 	d.plotterDataPacket = make([]byte, 0, plotter.PLOTTER_PACKET_LEN)
	// 	d.plotterDataPacket = append(d.plotterDataPacket, plotter.SerializeHeader()...)
	// }
	// d.SerializeRotations()
	// if d.plotterDataCounter < d.ploterDataPerPacket {
	// 	return false
	// }
	// if d.plotterUdpConn != nil {
	// 	copy(d.plotterSendBuffer, d.plotterDataPacket)
	// 	go func() {
	// 		d.plotterUdpConn.Write(d.plotterSendBuffer)
	// 	}()
	// }
	// d.plotterDataCounter = 0
	return true
}

func (d *droneApp) SerializeRotations() {
	// d.plotterDataPacket = append(
	// 	d.plotterDataPacket,
	// 	plotter.SerializeDroneData(
	// 		time.Since(d.startTime),
	// 		d.rotations,
	// 		d.accRotations,
	// 		d.gyroRotations,
	// 		0,
	// 	)...)
	// d.plotterDataCounter++
}
