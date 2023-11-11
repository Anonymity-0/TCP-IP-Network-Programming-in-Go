package main

import (
	"log"
	"net"
	"os"
)

func main() {

	//检查参数
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>\n", os.Args[0])
	}
	//设置监听端口和信息
	message := "Hello world!"
	port := os.Args[1]

	//当你调用 net.Listen("tcp", ":"+port) 时，Go 会执行以下操作：
	//1.创建一个 TCP 套接字。
	//2.将套接字绑定到指定的 IP 地址和端口号。
	//3.将套接字转为可接收连接状态
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("socket() error: %v", err)
	}
	//延迟关闭监听
	defer listener.Close()
	//循环监听
	for {
		//4.Accept()接收连接
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept() error: %v", err)
		}
		//5.write()发送信息
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Fatalf("write() error: %v", err)

		}

		conn.Close()
	}

}
