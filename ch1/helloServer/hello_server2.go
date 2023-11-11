package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
	}
	port, _ := strconv.Atoi(os.Args[1])

	//1.创建一个 TCP 套接字。
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Printf("socket() error: %v", err)
	}
	//2.将套接字绑定到指定的 IP 地址和端口号。
	// 创建一个 SockaddrInet4 结构体实例，并设置端口号
	servAddr := &syscall.SockaddrInet4{Port: port}
	// 将 IP 地址设置为 0.0.0.0，这样服务器可以接受来自任何客户端的连接
	copy(servAddr.Addr[:], []byte{0, 0, 0, 0})

	err = syscall.Bind(fd, servAddr)
	if err != nil {
		fmt.Printf("bind() error: %v", err)
	}
	//3.将套接字转为可接收连接状态
	err = syscall.Listen(fd, 10)
	if err != nil {
		fmt.Printf("listen() error: %v", err)
	}
	//4.Accept()接收连接
	connFd, _, err := syscall.Accept(fd)
	if err != nil {
		fmt.Printf("accept() error: %v", err)
	}
	//5.write()发送信息
	message := "Hello world!"
	_, err = syscall.Write(connFd, []byte(message))
	if err != nil {
		fmt.Printf("write() error: %v", err)
	}
	syscall.Close(connFd)
	syscall.Close(fd)

}
