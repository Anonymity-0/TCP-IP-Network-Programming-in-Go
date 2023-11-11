package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <IP> <port>\n", os.Args[0])
	}
	ip := os.Args[1]
	portStr := os.Args[2]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s\n", portStr)
	}

	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	// 设置服务器的地址和端口
	servAddr := &syscall.SockaddrInet4{Port: port}
	copy(servAddr.Addr[:], net.ParseIP(ip).To4())

	// 使用 syscall.Connect 函数将套接字连接到服务器
	err = syscall.Connect(fd, servAddr)
	if err != nil {
		log.Fatalf("connect() error: %v", err)
	}
	// 使用 syscall.Read 接收信息
	buf := make([]byte, 1024)
	message := make([]byte, 0)
	for {
		n, err := syscall.Read(fd, buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Read error: %v", err)
			}
			break
		}
		if n == 0 {
			break
		}
		message = append(message, buf[:n]...)
	}
	fmt.Printf("Message from server: %s\n", message)
	fmt.Printf("Function read() read %d bytes\n", len(message))
	syscall.Close(fd)
	syscall.SockaddrInet4

}
