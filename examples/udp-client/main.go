package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	const DATA_PER_SECOND = 100
	log.SetFlags(log.Lmicroseconds)
	udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:8000")
	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		log.Println(err)
		return
	}

	lastPrint := time.Now()

	for {
		data := genData()
		_, err = conn.Write(data)
		if time.Since(lastPrint) > time.Second/4 {
			fmt.Printf("sent %3d, Data/Second: %d\n", data[0], DATA_PER_SECOND)
			lastPrint = time.Now()
		}
		time.Sleep(time.Second / DATA_PER_SECOND)
	}
}

var index byte = 0

func genData() []byte {
	const SIZE = 6001
	data := make([]byte, SIZE)
	data[0] = index
	data[SIZE-1] = index
	if index == 255 {
		index = 0
	} else {
		index++
	}
	return data
}
