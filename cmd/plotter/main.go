package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Please provide host:port")
		os.Exit(1)
	}

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", "192.168.1.101:8009")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read from UDP listener in endless loop
	for {
		buf := make([]byte, 8912)
		n, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(n)

		// Write back the message over UPD
		// conn.WriteToUDP([]byte("Hello UDP Client\n"), addr)
	}
}
