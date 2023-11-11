package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <IP> <port>\n", os.Args[0])
	}

	ip := os.Args[1]
	port := os.Args[2]
	address := fmt.Sprintf("%s:%s", ip, port)
	//1.创建一个 TCP 套接字。
	//2.调用 Connect() 连接到指定的 IP 地址和端口号。
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	//net.Conn 对象来读取服务器发送的消息
	message, err := io.ReadAll(conn)
	if err != nil {
		log.Fatalf("Failed to read from server: %v", err)
	}

	fmt.Printf("Message from server: %s\n", message)
}
