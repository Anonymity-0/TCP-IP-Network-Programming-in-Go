package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const INADDR_NONE = 0xffffffff

// inet_addr converts an IPv4 address in dot-decimal notation into a 32-bit integer in network byte order.
// If the IP address is invalid, it returns INADDR_NONE (0xffffffff).
func inet_addr(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0xffffffff
	}
	ip = ip.To4()
	if ip == nil {
		return 0xffffffff
	}
	return binary.BigEndian.Uint32(ip)

}
func main() {
	addr1 := "1.2.3.4"
	addr2 := "1.2.3.256"

	conv_addr := inet_addr(addr1)

	if conv_addr == INADDR_NONE {
		fmt.Println("Error occured!")
	} else {
		fmt.Printf("Network ordered integer addr: %#x\n", conv_addr)
	}
	conv_addr = inet_addr(addr2)
	if conv_addr == INADDR_NONE {
		fmt.Println("Error occured!")
	} else {
		fmt.Printf("Network ordered integer addr: %#x\n", conv_addr)
	}
}
