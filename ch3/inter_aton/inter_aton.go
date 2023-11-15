package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

// inet_aton converts an IPv4 address in dot-decimal notation into a 32-bit integer in network byte order.
// If the IP address is valid, it stores the result in the given *uint32 and returns true.
// If the IP address is invalid, it returns false.
func inet_aton(ipStr string, ip *uint32) bool {
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return false
	}
	parsedIP = parsedIP.To4()
	if parsedIP == nil {
		return false
	}
	*ip = binary.BigEndian.Uint32(parsedIP)
	return true
}

func main() {
	addr := "127.232.124.79"
	var ip uint32

	if !inet_aton(addr, &ip) {
		log.Fatalln("Conversion error")
	} else {
		fmt.Printf("Network ordered integer addr: %#x\n", ip)
	}
}
