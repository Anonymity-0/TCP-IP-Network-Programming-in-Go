package main

import (
	"fmt"
	"net"
)

func main1() {
	var addr1, addr2 uint32 = 0x1020304, 0x1010101

	ip1 := net.IPv4(byte(addr1>>24), byte(addr1>>16), byte(addr1>>8), byte(addr1))
	ip2 := net.IPv4(byte(addr2>>24), byte(addr2>>16), byte(addr2>>8), byte(addr2))

	fmt.Printf("Dotted-Decimal notation1: %s \n", ip1.String())
	fmt.Printf("Dotted-Decimal notation2: %s \n", ip2.String())
}
