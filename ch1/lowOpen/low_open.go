package main

import (
	"log"
	"os"
)

func main() {

	buf := []byte("Let's go!\n")
	//打开文件，如果不存在则创建，如果存在则清空，权限为 0644
	f, err := os.OpenFile("data.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("open() error: %v", err)
	}
	defer f.Close()
	//输出文件描述符
	log.Printf("file descriptor: %d \n", f.Fd())
	//写入文件
	_, err = f.Write(buf)
	if err != nil {
		log.Fatalf("write() error: %v", err)
	}

}
