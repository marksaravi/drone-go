package main

import (
	"fmt"
	"net"
)

func createUdpConnection(appConfig ApplicationConfig) (
	udpCon *net.PacketConn,
	udpAddr *net.UDPAddr,
	udpEnabled bool) {

	udpAddr = nil
	udpCon = nil
	udpEnabled = false

	// if !appConfig.UDP.enabled {
	// 	fmt.Println("UDP is not enabled")
	// 	return
	// }
	con, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return
	}
	address, err := net.ResolveUDPAddr("udp", appConfig.UDP)
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return
	}
	udpCon = &con
	udpAddr = address
	udpEnabled = true
	return
}
