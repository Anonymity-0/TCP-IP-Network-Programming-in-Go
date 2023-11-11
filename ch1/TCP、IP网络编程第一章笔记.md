# 1.1 理解网络和套接字

## 套接字
套接字（socket）是计算机网络中进程间通信的一种机制，它将进程间通信抽象为一个端点，该端点由一个IP地址和一个端口号来标识。
套接字是网络数据传输用的软件设备。网络编程又称为套接字编程。


### 编写" Hello world! " 服 务 器 端
网络编程中接受连接请求的套接字创建过程可整理如下。
1. 第一步:调用s o c k e t 函数创建套接字。
2. 第二步:调用b i n d 函数分配E地址和端口号。
3. 第三步:调用l i s t e n 函数转为可接收请求状态。
4. 第四步:调用a c c e p t 函数受理连接请求。

hello_server.go
``` go
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


```


### client
创建套接字,但此时套接字并不马上分为服务器端和客户端。如果紧接着调用bind和listen函数,将成为**服务器端套接字**;如果调用connect函数将成为客户端套接字。

hello_client.go
```go
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


```

### 运行

- 编译运行服务器
```shell
go build -o ch1/hello_server ch1/hello_server.go
./ch1/hello_server 8080
```
- 编译运行客户端

```shell
go build -o ch1/hello_client ch1/hello_client/hello_client.go


./ch1/hello_client 127.0.0.1 8080      
Message from server: Hello world!
```


# 1.2 文件操作

## 打开文件
首先介绍打开文件以读写数据的函数。 调用此函数时需传递两个参数:第一个参数是打开的目标文件名及路径信息,第二个参数是文件打开模式(文件特性信息)。


``` c
#include <sys/types.h> 
#include <sys/stat.h> 
#include <fcntl.h>
int open(const char * path ，int flag)
//成功时返回文件描述符事失败时返回1
//path 文件名字符串地址
//flag 文件打开模式
```


```go
func Open(name string) (*File, error) {
	return OpenFile(name, O_RDONLY, 0)
}

// OpenFile 是一个通用的文件打开函数，它接受三个参数：文件名、打开文件的标志和文件权限。
// 它返回一个 *os.File 对象和一个 error 对象。如果打开文件失败，它会返回一个非 nil 的 error 对象。
func OpenFile(name string, flag int, perm FileMode) (*File, error) {
    // 记录打开文件的操作
    testlog.Open(name)
    // 使用 openFileNolog 函数打开文件，这个函数不记录日志
    f, err := openFileNolog(name, flag, perm)
    // 如果打开文件失败，返回错误
    if err != nil {
        return nil, err
    }
    // 检查打开文件的标志是否包含 O_APPEND，如果包含，设置 f.appendMode 为 true
    f.appendMode = flag&O_APPEND != 0

    // 返回打开的文件
    return f, nil
}
```

## 关闭文件

```c
#include <unistd.h> 
int close(int fd); 
	//成功时返回Q,失败时返回-1。
//fd:需要关闭的文件或套接字的文件描述符
```


```go
// Close closes the File, rendering it unusable for I/O.
// On files that support SetDeadline, any pending I/O operations will
// be canceled and return immediately with an ErrClosed error.
// Close will return an error if it has already been called.
func (f *File) Close() error {
	if f == nil {
		return ErrInvalid
	}
	return f.file.close()
}

```

## 将数据写入文件

```c
#include <unistd.h> 
ssize_t write(int fd, const void * buf, size_t nbytes); 
//成功时返回写入的字节数,失败时返回-1。
//fd:数据传输对象的文件描述符
//buf 保存数据的缓冲地址
//nbytes：要传输的字节数
```


```go
// write writes len(b) bytes to the File.
// It returns the number of bytes written and an error, if any.
// write 是 File 结构体的一个方法，它接受一个字节切片 b 作为参数，
// 尝试将这个字节切片写入到文件中。
func (f *File) write(b []byte) (n int, err error) {
    // f.pfd.Write(b) 调用 pfd（代表平台依赖的文件描述符）的 Write 方法，
    // 尝试将 b 写入到文件。这个方法返回写入的字节数和一个错误（如果有的话）。
    n, err = f.pfd.Write(b)
    // runtime.KeepAlive(f) 是一个用于防止 f 被垃圾回收的调用。
    // 在某些情况下，如果 f 在 f.pfd.Write(b) 调用之后没有被再次使用，
    // Go 的垃圾回收器可能会在 Write 调用还在进行时就回收 f。
    // runtime.KeepAlive(f) 确保 f 在 Write 调用完成之前不会被垃圾回收。
    runtime.KeepAlive(f)
    // 返回写入的字节数和错误（如果有的话）
    return n, err
}
```


### 代码示例改写

low_open.go

```go
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
```

**运行代码**

```shell
go build -o ch1/low_open ch1/lowOpen/low_open.go 
./ch1/low_open 
file descriptor: 3 
```

## 读取数据

```c
#include <unistd.h> 
ssize_t read(int fdJ void * buf, size_t nbytes); 
//'成功时返回接收的字节数(但遇到文件结尾则返回θ),失败时返回10
//fd 显示数据接收对象的文件描述符。
//buf 要保存接收数据的缓冲地址值。
//nbytes 要接收数据的最大字节数。

```

### 示例代码改写

low_read.go 
```go
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

```

**运行代码**

```shell
go build -o ch1/low_read  ch1/lowRead/low_read.go 
 
./ch1/low_read                                   
2023/11/10 16:45:54 file descriptor: 3 
2023/11/10 16:45:54 file data: Let's go!

```

## 文件描述符与套接字


原文c代码

```c
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/socket.h>

int main()
{
    int fd1, fd2, fd3;
    //创建一个文件和两个套接字
    fd1 = socket(PF_INET, SOCK_STREAM, 0);
    fd2 = open("test.dat", O_CREAT | O_WRONLY | O_TRUNC);
    fd3 = socket(PF_INET, SOCK_DGRAM, 0);
    //输出之前创建的文件描述符的整数值
    printf("file descriptor 1: %d\n", fd1);
    printf("file descriptor 2: %d\n", fd2);
    printf("file descriptor 3: %d\n", fd3);

    close(fd1);
    close(fd2);
    close(fd3);
    return 0;
}

```

改写后
fd_seri.go
```go
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

```

运行

```shell
go build -o ch1/fd_seri  ch1/fdSeri/fd_seri.go
./ch1/fd_seri 
file descriptor 1: 3
file descriptor 2: 4
file descriptor 3: 5

```

# 1.3 基于Windows平台的实现
略

# 1.4 基于Windows的套接字相关函数及示例


# 1.5 习题

1.  套接字在网络编程中的作用是什么?为何称它为套接字?
	
	套接字（socket）是网络编程中的抽象概念，它提供了一种机制，使得不同计算机之间可以进行通信和数据交换。套接字可以看作是网络通信的端点，它包含了通信所需的各种信息，如IP地址、端口号、协议等。套接字的名称来源于插座（socket），类比插座连接电器，套接字连接网络。通过套接字，计算机可以在网络上进行数据的发送和接收，实现网络通信的功能。

2. 在服务器端创建套接字后,会依次调用listen函数和accept 函数。请比较并说明二者作用。
	
	listen函数用于将套接字标记为被动套接字，即用于接受客户端的连接请求。它告诉操作系统该套接字将用于接受传入的连接，而不是发起连接。在调用listen函数后，套接字将进入监听状态，等待客户端的连接请求。
	accept函数用于从处于监听状态的套接字中接受一个连接。当客户端尝试连接到服务器时，服务器调用accept函数来接受这个连接，并创建一个新的套接字来与客户端进行通信。这个新的套接字可以用于与该客户端进行数据交换，而原始的监听套接字则继续等待其他客户端的连接请求。accept函数的返回值是一个新的套接字，通过它可以进行与客户端的通信。

3. Linux中,对套接字数据进行I/O时可以直接使用文件I/O 相关函数;而在Windows中则不可以。原因为何?
	这是因为在Linux中，套接字被视为一种文件描述符，因此可以使用文件I/O相关函数（如read和write）来进行I/O操作。而在Windows中，套接字和文件描述符是不同的概念，Windows采用了不同的I/O模型，因此不能直接使用文件I/O相关函数来对套接字数据进行I/O操作。在Windows中，需要使用特定的套接字I/O函数（如recv和send）来进行套接字数据的读写操作。

4. 创建套接字后一般会给它分配地址,为什么?为了完成地址分配需要调用哪个函数
	创建套接字后需要给它分配地址，这是为了让其他主机能够找到并与该套接字进行通信。在网络编程中，这个地址通常是IP地址和端口号的组合。
	为了完成地址分配，需要调用bind函数。bind函数将一个本地地址（IP地址和端口号）分配给套接字，使得其他主机可以通过这个地址与该套接字进行通信。

5. Linux中的文件描述符与Windows的句柄实际上非常类似。请以套接字为对象说明它们的含义。
	文件描述符和Windows的句柄在套接字的上下文中具有类似的含义。它们都是用来标识和引用套接字的抽象概念。
	在Linux中，套接字也被视为一种文件描述符，因此可以使用类似于文件I/O的操作来进行套接字的读写等操作。
	在Windows中，套接字使用句柄来进行引用和操作，句柄是一种抽象的引用类型，可以用来标识和操作套接字。
	因此，无论是文件描述符还是句柄，它们都是用来引用和操作套接字这种抽象对象的标识符。

6. 底层文件I/O函数与ANSI 标准定义的文件I/O函数之间有何区别?
	底层文件I/O函数是直接调用操作系统提供的文件操作接口，如open、read、write等，它们提供了对文件的低级别访问，可以更加灵活地控制文件的读写操作。
	而ANSI标准定义的文件I/O函数则是标准C库中提供的一组文件操作函数，如fopen、fread、fwrite等，它们提供了更加抽象和便捷的文件操作接口，使得跨平台开发更加方便，并且提供了一些缓冲和错误处理的功能。
	因此，底层文件I/O函数更加接近操作系统提供的文件操作接口，而ANSI标准定义的文件I/O函数则提供了更加便捷和跨平台的文件操作接口。

7. 参考本书给出的示例low_open.c 和low_read.c ,分别利用底层文件I/O 和ANSI标准I/O 编写文件复制程序。可任意指定复制程序的使用方法。