package main

import (
	"encoding/binary"
	"fmt"
)

func htons(n uint16) uint16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return binary.BigEndian.Uint16(b)
}

func htonl(n uint32) uint32 {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return binary.BigEndian.Uint32(b)
}

func main() {
	hostPort := uint16(0x1234)
	hostAddr := uint32(0x12345678)

	netPort := htons(hostPort)
	netAddr := htonl(hostAddr)

	fmt.Printf("Host ordered port: %#x \n", hostPort)
	fmt.Printf("Network ordered port: %#x \n", netPort)
	fmt.Printf("Host ordered address: %#x \n", hostAddr)
	fmt.Printf("Network ordered address: %#x \n", netAddr)
}
