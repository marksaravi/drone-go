package main

import (
	"log"
	"net"
	"time"
)

const BUFFER_SIZE = 8192

func main() {
	log.SetFlags(log.Lmicroseconds)
	udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:8000")
	// Start listening for UDP packages on the given address
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Println(err)
		return
	}

	buffer := make([]byte, BUFFER_SIZE)
	log.Println("Waiting for data...")
	lastPrint := time.Now()
	dataPerSecond := 0
	for {
		nBytes, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
		} else {
			dataPerSecond++
			if time.Since(lastPrint) >= time.Second {
				lastPrint = time.Now()
				log.Printf("Packet Size: %4d, Data Per Second: %3d\n", nBytes, dataPerSecond)
				dataPerSecond = 0
			}

		}
	}

}
