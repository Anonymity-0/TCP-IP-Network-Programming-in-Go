package main

import (
	"fmt"

	"os"
	"syscall"
)

func main() {
	// 使用 syscall.Socket 函数创建一个 TCP 套接字，返回的是文件描述符 fd1
	fd1, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	// 使用 os.OpenFile 函数创建一个文件 "test.dat"，返回的是 *os.File 类型的 fd2
	fd2, _ := os.OpenFile("test.dat", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	// 使用 syscall.Socket 函数创建一个 UDP 套接字，返回的是文件描述符 fd3
	fd3, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	// 使用 fmt.Printf 函数打印出三个文件描述符的值
	// 注意，对于 *os.File 类型的 fd2，我们使用了 Fd 方法来获取其底层的文件描述符
	fmt.Printf("file descriptor 1: %d\n", fd1)
	fmt.Printf("file descriptor 2: %d\n", fd2.Fd())
	fmt.Printf("file descriptor 3: %d\n", fd3)

	// 使用 Close 方法和 syscall.Close 函数关闭了这三个文件描述符，以释放系统资源
	fd2.Close()
	syscall.Close(fd1)
	syscall.Close(fd3)
}
