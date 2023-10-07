package main

import "github.com/marksaravi/drone-go/apps/plotter"

func main() {
	p := plotter.NewPlotter(plotter.PlotterSettings{
		UDPServerAddress:  ":4014",
		HTTPServerAddress: ":3000",
	})
	p.StartPlotter()
	// // Resolve the string address to a UDP address
	// udpAddr, err := net.ResolveUDPAddr("udp", "192.168.1.101:8009")

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// // Start listening for UDP packages on the given address
	// conn, err := net.ListenUDP("udp", udpAddr)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// // Read from UDP listener in endless loop
	// for {
	// 	buf := make([]byte, 8912)
	// 	n, _, err := conn.ReadFromUDP(buf)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	fmt.Println(n)

	// 	// Write back the message over UPD
	// 	// conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
	// }
}
