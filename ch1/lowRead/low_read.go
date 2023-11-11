package main

import (
	"log"
	"os"
)

func main() {
	//只读打开文件
	f, err := os.Open("data.txt")
	if err != nil {
		log.Fatalf("open() error: %v", err)
	}
	defer f.Close()

	//输出文件描述符
	log.Printf("file descriptor: %d \n", f.Fd())

	//读取文件
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		log.Fatalf("read() error: %v", err)
	}
	log.Printf("file data: %s", buf[:n])
}
