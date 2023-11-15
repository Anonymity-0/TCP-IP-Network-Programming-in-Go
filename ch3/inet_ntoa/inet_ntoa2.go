package main

import (
	"fmt"
	"net"
)

// inet_ntoa converts a 32-bit integer in network byte order into a dotted-decimal IP address.
func inet_ntoa(ipInt uint32) string {
	ipBytes := make([]byte, 4)
	ipBytes[0] = byte(ipInt >> 24)
	ipBytes[1] = byte(ipInt >> 16)
	ipBytes[2] = byte(ipInt >> 8)
	ipBytes[3] = byte(ipInt)
	return net.IP(ipBytes).String()
}

func main() {
	var addr1, addr2 uint32 = 0x1020304, 0x1010101

	fmt.Printf("Dotted-Decimal notation1: %s \n", inet_ntoa(addr1))
	fmt.Printf("Dotted-Decimal notation2: %s \n", inet_ntoa(addr2))
}
