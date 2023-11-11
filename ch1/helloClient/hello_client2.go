package main

import (
	"fmt"
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
	} else {
		fmt.Println("connect success")
	}
	// 使用 syscall.Read 接收信息
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		log.Fatalf("read() error: %v", err)
	}
	// 打印接收到的信息
	fmt.Printf("Message from server: %s\n", buf[:n])
	fmt.Printf("Function read() read %d bytes\n", n)
	syscall.Close(fd)

}
