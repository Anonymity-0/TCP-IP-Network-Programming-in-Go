
# 2.1 套接字协议及其数据传输特性

## 协议
如果相隔很远的两人想展开对话,必须先决定对话方式。如果一方使用电话,那么另一方也只能使用电话,而不是书信。可以说,电话就是两人对话的协议。协议是对话中使用的通信规则, 把上述概念拓展到计算机领域可整理为“**计算机间对话必备通信规则**”。

## 创建套接字

```c
#include <sys/socket.h> 
int socket (int domain, int type, int protocol); 
//成功时返回文件描述符,失败时返回-1。
//domain 套接字中使用的协议族(Protocol Family)信息。
//type 套接字数据传输类型信息。
//protocol 计算机间通信中使用的协议信息。
```

在 Go 语言的 `syscall` 包中，`Socket` 函数是对底层 `socket` 系统调用的封装。这样做的目的是为了提供一个更加 Go 风格（例如错误处理）的接口，同时隐藏一些底层细节。

`Socket` 函数内部调用了 `socket` 函数。`socket` 函数直接执行了系统调用，并返回了原始的结果，包括一个文件描述符和一个错误号。然后 `Socket` 函数将这些原始结果转换为 Go 风格的结果：如果系统调用成功，它返回一个文件描述符和一个 `nil` 错误；如果系统调用失败，它返回一个 `-1` 文件描述符和一个非 `nil` 错误。

这样做的好处是，对于大多数 Go 程序员来说，他们只需要关心 `Socket` 函数，而不需要了解底层的 `socket` 系统调用和错误处理。

```go
// Socket 函数创建一个新的套接字，并返回其文件描述符和可能的错误。
// domain 参数指定了套接字的协议族（例如，AF_INET 代表 IPv4，AF_INET6 代表 IPv6）。
// typ 参数指定了套接字的类型（例如，SOCK_STREAM 代表 TCP，SOCK_DGRAM 代表 UDP）。
// proto 参数指定了套接字使用的协议（例如，IPPROTO_TCP 代表 TCP，IPPROTO_UDP 代表 UDP）。
func Socket(domain, typ, proto int) (fd int, err error) {
    // 如果 domain 是 AF_INET6（即，我们正在尝试创建一个 IPv6 套接字），
    // 但是 SocketDisableIPv6 为 true（即，我们禁用了 IPv6），则返回错误 EAFNOSUPPORT。
    if domain == AF_INET6 && SocketDisableIPv6 {
        return -1, EAFNOSUPPORT
    }
    // 调用底层的 socket 函数创建套接字。
    fd, err = socket(domain, typ, proto)
    return
}




// socket 函数创建一个新的套接字，并返回其文件描述符和可能的错误。
// domain 参数指定了套接字的协议族（例如，AF_INET 代表 IPv4，AF_INET6 代表 IPv6）。
// typ 参数指定了套接字的类型（例如，SOCK_STREAM 代表 TCP，SOCK_DGRAM 代表 UDP）。
// proto 参数指定了套接字使用的协议（例如，IPPROTO_TCP 代表 TCP，IPPROTO_UDP 代表 UDP）。
// rawSyscall 函数执行一个底层的系统调用，其参数是系统调用的编号和参数。
// 如果系统调用失败，它返回一个非零的错误号 e1，我们将其转换为 Go 的 error 类型并返回。
func socket(domain int, typ int, proto int) (fd int, err error) {
    r0, _, e1 := rawSyscall(abi.FuncPCABI0(libc_socket_trampoline), uintptr(domain), uintptr(typ), uintptr(proto))
    fd = int(r0)
    if e1 != 0 {
        err = errnoErr(e1)
    }
    return
}

```


## 协议族(Protocol Family)
通过socket函数的第一个参数传递套接字中使用的协议分类信息。此协议分类信息称为协议族。
原文给的c语言头文件中的分类
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.3ptuslo70nq.png)
在 Go 语言的 `syscall` 包中，`domain` 参数用于指定套接字的协议族。以下是一些常见的 `domain` 值：

|值|描述|
|---|---|
|`AF_INET`|IPv4 网络协议|
|`AF_INET6`|IPv6 网络协议|
|`AF_UNIX`|UNIX 域套接字|
|`AF_NETLINK`|内核用户接口设备|
|`AF_PACKET`|低级包接口|
|`AF_CAN`|Controller Area Network|
|`AF_BLUETOOTH`|蓝牙设备|

请注意，不是所有的 `domain` 值都在所有的平台上可用。具体可用的 `domain` 值取决于你的操作系统和平台。

## 套接字类型(Type)
套接字类型指的是套接字的数据传输方式,通过socket函数的第二个参数传递,只有这样才能决定创建的套接字的数据传输方式。这种说法可能会使各位感到疑惑。已通过第一个参数传递了协议族信息,还要决定数据传输方式?问题就在于,决定了协议族并不能同时决定数据传输方式,换言之,socket数第一个参数PF_INET协议族中也存在多种数据传输方式。

### 面向连接的套接字(SOCK_STREAM)
如果向socket函数的第二个参数传递SOCK_STREAM,将创建面向连接的套接字。面向连接的套接字到底具有哪些特点呢?

- 套接字连接必须一一对应
- 传输过程中数据不会消失
- 按序传输数据
- 传输的数据不存在数据边界


收发数据的套接字内部有缓冲(buffer),简言之就是字节数组。通过套接字传输的数据将保存到该数组。因此,收到数据并不意味着马上调用read函数。只要不超过数组容量,则有可能在数据填充满缓冲后通过1次read函数调用读取全部,也有可能分成多次read函数调用进行读取。也就是说,在面向连接的套接字中,read函数和write函数的调用次数并无太大意义。所以说面向连接的套接字不存在数据边界。


### 套接字缓冲已满是否意味着数据丢失
之前讲过,为了接收数据,套接字内部有一个由字节数组构成的缓冲。如果这个缓冲被接收的数据填满会发生什么事情?之后传递的数据是否会丢失? 
首先调用read函数从缓冲读取部分数据,因此,缓冲并不总是满的。但如果read函数读取速度比接收数据的速度慢,则缓冲有可能被填满。此时套接字无法再接收数据, 但即使这样也不会发生数据丢失,因为**传输端套接字将停止传输**。也就是说,**面向连接的套接字会根据接收端的状态传输数据,如果传输出错还会提供重传服务**。因此,**面向连接的套接字除特殊情况外不会发生数据丢失**。


### 面向消息的套接字(SOCK_DGRAM)
如果向socket函数的第二个参数传递SOCK_DGRAM,则将创建面向消息的套接字。面向消息的套接字可以比喻成高速移动的摩托车快递。
- 强调快速传输而非传输顺序
- 传输的数据可能丢失也可能损毁
- 传输的数据有数据边界
- 限制每次传输的数据大小

众所周知,快递行业的速度就是生命。用摩托车发往同一目的地的2件包裹无需保证顺序, 只要以最快速度交给客户即可。这种方式存在损坏或丢失的风险,而且包裹大小有一定限制。因此,若要传递大量包裹,则需分批发送。另外,**如果用2辆摩托车分别发送2件包裹,则接收者也需要分2次接收**。这种特性就是“传输的数据具有**数据边界**”。

面向消息的套接字比面向连接的套接字具有更快的传输速度,但无法避免数据丢失或损毁。另外,每次传输的数据大小具有一定限制,并存在数据边界。存在数据边界意味着接收数据的次数应和传输次数相同。面向消息的套接字特性总结如下: “**不可靠的、不按序传递的、以数据的高速传输沟目的的套接字**"


### 最终选择的协议

下面讲解socket函数的第三个参数,该参数决定最终采用的协议。
前面已经通过socket函数的前两个参数传递了协议族信息和套接字数据传输方式,这些信息还不足以决定采用的协议吗?为什么还需要传递第3个参数呢? 
传递前两个参数即可创建所需套接字。所以大部分情况下可以向第三个参数传递**0**,除非遇到以下这种情况: “**同一协议族中存在多个数据传输方式相同的协议**”数据传输方式相同,但协议不同。此时需要通过第三个参数具体指定协议信息。


### TCP套接字示例

tcp_server.go (系统函数调用版)

```go
package main

import (
	"fmt"
	"log"
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
		log.Fatalf("bind() error: %v", err)
	}
	//3.将套接字转为可接收连接状态
	err = syscall.Listen(fd, 10)
	if err != nil {
		log.Fatalf("listen() error: %v", err)
	}
	//4.Accept()接收连接
	connFd, _, err := syscall.Accept(fd)
	if err != nil {
		log.Fatalf("accept() error: %v", err)
	}
	//5.write()发送信息
	message := "Hello world!"
	_, err = syscall.Write(connFd, []byte(message))
	if err != nil {
		log.Fatalf("write() error: %v", err)
	}
	syscall.Close(connFd)
	syscall.Close(fd)

}


```

tcp_clinet.go 更改read函数调用方式，在客户端中分多次调用read函数以接收服务器端发送的全部数据,以验证tcp传输的数据不存在数据边界。

```go
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

}


```
运行

```shell
go build -o ch2/tcpServer/tcp_server ch2/tcpServer/tcp_server.go 
go build -o ch2/tcpClient/tcp_client ch2/tcpClient/tcp_client.go 

./ch2/tcpServer/tcp_server 8888              
./ch2/tcpClient/tcp_client 127.0.0.1 8888

Message from server: Hello world!
Function read() read 12 bytes
```

# 2.2 Windows平台啊下的实现及验证

略


# 2.3 习题
1. 什么是协议?在收发数据中定义协议有何意义? 

2. 面向连接的TCP套接字传输特性有3点,请分别说明。
	- 无数据边界
	- 传输过程中数据不会消失
	- 按序传输
3. 下列哪些是面向消息的套接字的特性? 
	a. <font color="#ff0000">传输数据可能丢失</font>
	b. ~~没有数据边界(Boundary)~~ （面向连接）
	c. <font color="#ff0000">以快速传递为目标</font>
	d. ~~不限制每次传递数据的大小~~ （限制大小）
	<font color="#ff0000">e. 与面向连接的套接字不同,不存在连接的概念</font>
	 
4. 下列数据适合用哪类套接字传输?并给出原因。
	1. a. 演唱会现场直播的多媒体数据(面向消息) 
		因为面向消息的套接字以快速传递为目标，适合传输多媒体数据，即使传输数据可能丢失也不会影响整体效果。
	2. b. 某人压缩过的文本文件(面向连接) 
		因为面向连接的套接字可以保证数据的可靠传输，适合传输对数据完整性要求较高的文本文件。
	3. c. 网上银行用户与银行之间的数据传递(面向连接) 
		面向连接的套接字提供可靠的、按顺序传送的数据传输服务，适合对数据完整性和安全性要求较高的网上银行交易。
	
5. 何种类型的套接字不存在数据边界?这类套接字接收数据时需要注意什么? 
	面向连接的套接字不存在数据边界。面向连接的TCP套接字在接收数据时需要注意处理粘包和拆包的问题，确保按照应用层协议的要求正确解析和处理接收到的数据
6. tcp_server.c和tcp_client.c中需多次调用read函数读取服务器端调用I次write函数传递的字符串。更改程序,使服务器端多次调用(次数自拟)write函数传输数据,客户端调用1 次read函数进行读取。为达到这一目的,客户端需延迟调用read函数,因为客户端要等待服务器端传输所有数据。Windows和Linux都通过下列代码延迟read或recv函数的调用。for(1=0;i<3000;i++) printf("wait time %d \n", i); 让CPU执行多余任务以延迟代码运行的方式称为 “Busy Waiting”。使用得当即可推迟函数调用。