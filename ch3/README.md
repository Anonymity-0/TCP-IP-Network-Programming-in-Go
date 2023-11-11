
## 3.1 分配给套接字的IP地址与端口号

### 网络地址
略


## 用于区分套接字的端口号

IP用于区分计算机,只要有IP地址就能向目标主机传输数据,但仅凭这些无法传输给最终的应用程序。假设各位欣赏视频的同时在网上冲浪,这时至少需要1个接收视频数据的套接字和1 个接收网页信息的套接字。问题在于如何区分二者。简言之,传输到计算机的网络数据是发给播放器,还是发送给浏览器?

若想接收多台计算机发来的数据,则需要相应个数的套接字。那如何区分这些套接字呢? 
计算机中一般配有NIC(Network Interface Card,网络接口卡)数据传输设备。通过NIC向计算机内部传输数据时会用到IP。操作系统负责把传递到内部的数据适当分配给套接字,利用端口号。也就是说,通过NIC接收的数据内有端口号,操作系统正是参考此端口号把数据传输给相应端口的套接字
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.5q51v4z8i4n4.webp)

端口号就是在同一操作系统内为区分不同套接字而设置的,因此无法将1个端口号分配给不同套接字。
另外,端口号由16位构成 可分配的端口号范围是**0-65535**。但**0-1023**是**知名端口(Well-known PORT)**,一般分配给特定应用程序,所以应当分配此范围之外的值。另外,虽然端口号不能重复,**但TCP套接字和UDP套接字不会共用端口号,所以允许重复**。例如:如果某TCP 套接字使用9190号端口,则其他TCP套接字就无法使用该端口号,但UDP套接字可以使用。
总之,数据传输目标地址同时包含IP地址和端口号,只有这样,数据才会被传输到最终的目的应用程序(应用程序套接字)。

## 3.2
应用程序中使用的IP地址和端口号以结构体的形式给出了定义。本节将以IPv4为中心,围绕此结构体讨论目标地址的表示方法。

### 表示IPv4的结构体

填写地址信息时应以如下提问为线索进行

	口 问题1：“采用哪一种地址族？”
	口 答案1：“基于IPv4的地址族。”
	口问题2：“IP地址是多少？”
	口答案2：“211.204.214.76。”
	口 问题3：“端口号是多少？”
	口 答案3：“2048。”

#### C
C语言中IPv4结构体定义为如下形态

```c
struct sockaddr_in
{
	sa_family_t      sin_family; //地址族
	uint16_t         sin_port;   //16位TCP/UDP地址
	struct in_addr          sin_addr;  //32位ip地址
	char             sin_zero[8];    //不使用
};

```

`in_addr`定义如下，它用来存放32位IP地址

```c
struct in_addr
{
	In_addr_t    s_addr;  //32位IPv4地址
}

```
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.1a3tohqxb3ts.webp)

可以看到 `in_addr_t` 实际上是 `uint32_t`也就是无符号32位bit，那为什么需要额外定义这些数据类型呢? 如前所述,这是考虑到扩展性的结果。如果使用int32_t类型的数据,就能保证在任何时候都占用4字节,即使将来用64位表示int类型也是如此。
从之前介绍的代码也可看出,`sockaddr_in`结构体变量地址值将以如下方式传递给`bind`函数。

``` c
struct sockaddr_in serv_addr; 
...
if(bind(serv_sock,(struct sockaddr * ) &serv_addr, sizeof(Serv_addr))== -1) 
	error_handling("bind()error");
...
```
此处重要的是第二个参数的传递。实际上,bind函数的第二个参数期望得到`sockaddr`结构体变量地址值,包括地址族、端口号、IP地址等。（此处进行了强制类型转换,将`sockaddr_in`转成`sockaddr`）从下列代码也可看出,直接向sockaddr结构体填充这些信息会带来麻烦。

```c
struct sockaddr {
    sa_family_t char sin_family;//地址族(Address Family) 
    sa_data[14];// 地址信息
}
```
此结构体成员`sa_data`保存的地址信息中需包含IP地址和端口号,剩余部分应填充0,这也是`bind`函数要求的。而这对于包含地址信息来讲非常麻烦,继而就有了新的结构体`sockaddr_in`。若按照之前的讲解填写`sockaddr_in`结构体,则将生成符合bind函数要求的字节流。最后转换为`sockaddr`型的结构体变量,再传递给`bind`函数即可。


- `sin_family`
	每种协议族适用的地址族均不同。比如,IPv4使用4字节地址族,IPv6使用16字节地址族。
	可以参考表3-2保存`sin_family`地址信息。
	![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.43kn44zq1f9c.webp)
	
	sockaddr_in是保存IPv4地址信息的结构体。那为何还需要通过sin_family单独指定地址族信息呢?
	这与之前讲过的`sockaddr`结构体有关。结构体sockaddr并非只次IPv4设计, 这从保存地址信息的数组`sa_data`长度为14字节也可看出。因此,结构体sockaddr要求在sin_family 中指定地址族信息。为了与sockaddr保持一致,sockaddr_in结构体中也有地址族信息。

- `sin_port`
	该成员保存16位端口号,重点在于,它以网络字节序保存

- `sin_addr`
	该成员保存32位IP地址信息,且也以网络字节序保存。为理解好该成员,应同时观察结构体`in_addr`。但结构体`in_addr`声明为`uint32_t`,因此只需当作32位整数型即可。

- `sin_zero`
	无特殊含义。只是 使结构体`sockaddr_in`的大小与`sockadd`结构体保持一致而插入的成员。
	必需填充为0,否则无法得到想要的结果。
#### go
在 Go 语言中，类似于 C 语言中 `struct sockaddr_in` 的结构体是 `syscall.SockaddrInet4`（对于 IPv4）和 `syscall.SockaddrInet6`（对于 IPv6）。

在go中，就没有特别定义`in_addr_t`，直接使用\[4]byte。
```go
// SockaddrInet4 结构体用于表示一个 IPv4 网络地址。
type SockaddrInet4 struct {
    Port int       // Port 字段表示端口号
    Addr [4]byte   // Addr 字段表示 IPv4 地址，存储为 4 字节
    raw  RawSockaddrInet4  // raw 字段是内部使用的原始结构体，用于与系统调用交互
}
```


-  `Port`：一个 `int` 类型的值，表示网络地址的端口号。

- `Addr`：一个 `[4]byte` 类型的数组，表示 IPv4 地址。每个字节代表地址的一部分，例如，地址 "127.0.0.1" 会被表示为 `[127, 0, 0, 1]`。

- `raw`：一个 `RawSockaddrInet4` 类型的值，表示网络地址的底层表示。这个字段通常由系统调用使用，不应在常规 Go 代码中直接使用


**RawSockaddrInet4**

```go
// RawSockaddrInet4 结构体用于表示一个 IPv4 网络地址的底层结构。
type RawSockaddrInet4 struct {
    Len    uint8     // Len 字段表示此结构体的长度
    Family uint8     // Family 字段表示地址族，对于 IPv4，此值通常为 AF_INET
    Port   uint16    // Port 字段表示端口号
    Addr   [4]byte   // Addr 字段表示 IPv4 地址，存储为 4 字节
    Zero   [8]int8   // Zero 字段是填充字段，用于确保结构体的大小正确
}

```
Go语言也类似，这是go 的bind函数

```go
func Bind(fd int, sa Sockaddr) (err error) {
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return bind(fd, ptr, n)
}

```


在 Go 语言中，`syscall.Bind` 函数的第二个参数是 `Sockaddr` 类型，这是一个接口类型，它定义了一些方法，这些方法需要由任何实现该接口的类型来实现。

**`SockaddrInet4` 和 `SockaddrInet6` 结构体都实现了 `Sockaddr` 接口**，因此它们可以作为 `syscall.Bind` 函数的参数。

当你创建一个 `SockaddrInet4` 结构体并传递给 `syscall.Bind` 函数时，**Go 语言会自动将 `SockaddrInet4` 结构体转换为 `Sockaddr` 接口类型**，然后再传递给 `syscall.Bind` 函数。

所以，虽然 `syscall.Bind` 函数的参数类型是 `Sockaddr`，但你可以传递一个 `*SockaddrInet4` 或 `*SockaddrInet6` 结构体给它。


- `Len`：一个 `uint8` 类型的值，表示此结构体的长度。
- `Family`：一个 `uint8` 类型的值，表示地址族。对于 IPv4，此值通常为 `AF_INET`。
	
	`Family` 字段在 `RawSockaddrInet4` 结构体中用于指定地址族。对于 IPv4 地址，这个字段通常被设置为 `AF_INET`。
	虽然 `RawSockaddrInet4` 结构体通常用于表示 IPv4 地址，但 `Family` 字段仍然是必要的，因为它告诉操作系统如何解释这个结构体中的其他字段。例如，`Port` 和 `Addr` 字段的解释方式取决于 `Family` 字段的值。
	此外，`Family` 字段也可以帮助调试和错误检查。例如，如果你看到一个 `Family` 字段的值不是 `AF_INET`，但结构体是 `RawSockaddrInet4`，那么你就知道有些地方出错了。

- `Port`：一个 `uint16` 类型的值，表示网络地址的端口号。注意，这个值是网络字节序。

- `Addr`：一个 `[4]byte` 类型的数组，表示 IPv4 地址。每个字节代表地址的一部分，例如，地址 "127.0.0.1" 会被表示为 `[127, 0, 0, 1]`。

- `Zero`：一个 `[8]int8` 类型的数组，用于填充，以确保结构体的大小正确。这个字段通常不用于常规编程。（类似sin_zero）


`SockaddrInet4` 和 `RawSockaddrInet4` 两个结构体都包含 `Port` 和 `Addr` 字段，但它们的用途是不同的。

`SockaddrInet4` 是 Go 语言对网络地址的高级表示，它的 `Port` 和 `Addr` 字段类型分别为 `int` 和 `[4]byte`，这对于 Go 程序员来说更易于使用。

而 `RawSockaddrInet4` 是对系统调用级别的网络地址的低级表示，它的 `Port` 和 `Addr` 字段类型分别为 `uint16` 和 `[4]byte`，并且 `Port` 字段是网络字节序，这对于系统调用来说是必须的。

当你在 Go 代码中创建一个 `SockaddrInet4` 结构体并传递给如 `syscall.Bind` 这样的函数时，Go 语言会自动将 `SockaddrInet4` 结构体转换为 `RawSockaddrInet4` 结构体，然后再传递给底层的系统调用。这就是为什么 `SockaddrInet4` 结构体中包含一个 `RawSockaddrInet4` 字段的原因。

`SockaddrInet4`结构体包含了IPv4地址和端口信息，而`RawSockaddrInet4`结构体则是为了在底层网络编程中使用原始的套接字地址结构而定义的。这种设计可以让网络编程在不同层次上进行操作，同时保持灵活性和可扩展性。



## 3.3 网络字节序与地址变换
不同CPU中,4字节整数型值1在内存空间的保存方式是不同的。4字节整数型值1可用2进制表示如下。
	00000000 00000000 00000000 00000001
有些CPU以这种顺序保存到内存,另外一些CPU则以倒序保存。
	00000001 00000000 00000000 00000000 
若不考虑这些就收发数据则会发生问题,因为保存顺序的不同意味着对接收数据的解析顺序也不同。

### 字节序与网络字节序
CPU向内存保存数据的方式有2种,这意味着CPU解析数据的方式也分为2种。
- 大端序(Big Endian):高位字节存放到低位地址。
	![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.27xmkxdgjg74.webp)
- 小端序(Little Endian):高位字节存放到高位地址。
	![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.2gkyxou8dm4g.webp)
每种CPU的数据保存方式均不同。因此, 代表CPU数据保存方式的主机字节序(Host Byte Order)在不同CPU中也各不相同。目前主流的Intel系列CPU以小端序方式保存数据。接下来分析2台字节序不同的计算机之间数据传递过程中可能出现的问题
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.4fi8nqamn5kw.webp)
0x12和0x34构成的大端序系统值与0x34和0x12构成的小端序系统值相同。换言之,只有改变数据保存顺序才能被识别为同一值。图3-6中,大端序系统传输数据0x1234时未考虑字节序问题, 而直接以Ox12、0x34的顺序发送。结果接收端以小端序方式保存数据,因此小端序接收的数据变成0x3412,而非0x1234。正因如此,在通过网络传输数据时约定统一方式,这种约定称为**网络字节序(Network Byte Order)** 非常简单：**统一为大端序**。
即,先把数据数组转化成大端序格式再进行网络传输。因此,所有计算机接收数据时应识别该数据是网络字节序格式,小端序系统传输数据时应转化为大端序排列方式。


### 字节序转换
接下来介绍帮助转换字节序的函数。这是文中给的c语言转换函数
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.2mqwp5w3kkcg.webp)

在 Go 语言中，没有内置的 `htons`（Host TO Network Short）、`ntohs`（Network TO Host Short）、`htonl`（Host TO Network Long）和 `ntohl`（Network TO Host Long）函数。这些函数在 C 语言中用于在主机字节序和网络字节序之间转换数据。

但是，你可以使用 Go 语言的 `encoding/binary` 包来实现同样的功能。以下是如何在 Go 语言中实现这些函数的例子：

```go
package main

import (
    "encoding/binary"
    "fmt"
)

func htons(n uint16) uint16 {
    var b [2]byte
    binary.BigEndian.PutUint16(b[:], n)
    return binary.BigEndian.Uint16(b[:])
}

func ntohs(n uint16) uint16 {
    return htons(n) // 在 16 位无符号整数上，ntohs 和 htons 是相同的
}

func htonl(n uint32) uint32 {
    var b [4]byte
    binary.BigEndian.PutUint32(b[:], n)
    return binary.BigEndian.Uint32(b[:])
}

func ntohl(n uint32) uint32 {
    return htonl(n) // 在 32 位无符号整数上，ntohl 和 htonl 是相同的
}

func main() {
    fmt.Println(htons(12345)) // 输出：12345
    fmt.Println(ntohs(12345)) // 输出：12345
    fmt.Println(htonl(12345)) // 输出：12345
    fmt.Println(ntohl(12345)) // 输出：12345
}

```


endian_conv.go

```go
package main

import (
    "encoding/binary"
    "fmt"
)

// htons 函数接收一个 uint16 类型的主机字节序值，
// 将其转换为网络字节序，然后返回转换后的值。
func htons(n uint16) uint16 {
    b := make([]byte, 2) // 创建一个长度为 2 的字节数组
    binary.BigEndian.PutUint16(b, n) // 将 n 的值以大端字节序放入字节数组
    return binary.BigEndian.Uint16(b) // 从字节数组中读取并返回大端字节序的值
}

// htonl 函数接收一个 uint32 类型的主机字节序值，
// 将其转换为网络字节序，然后返回转换后的值。
func htonl(n uint32) uint32 {
    b := make([]byte, 4) // 创建一个长度为 4 的字节数组
    binary.BigEndian.PutUint32(b, n) // 将 n 的值以大端字节序放入字节数组
    return binary.BigEndian.Uint32(b) // 从字节数组中读取并返回大端字节序的值
}

func main() {
    hostPort := uint16(0x1234) // 定义一个主机字节序的端口值
    hostAddr := uint32(0x12345678) // 定义一个主机字节序的地址值

    netPort := htons(hostPort) // 将主机字节序的端口值转换为网络字节序
    netAddr := htonl(hostAddr) // 将主机字节序的地址值转换为网络字节序

```


## 3.4 网络地址的初始化与分配

### 将字符串信息转换为网络字节序的整数型
sockaddr_in中保存地址信息的成员为32位整数型。因此,为了分配IP地址,需要将其表示为32位整数型数据。这对于只熟悉字符串信息的我们来说实非易事。

对于IP地址的表示,我们熟悉的是点分十进制表示法(Dotted Decimal Notation),而非整数型数据表示法。幸运的是,有个函数会帮我们将字符串形式的IP地址转换成32位整数型数据。此函数在转换类型的同时进行网络字节序转换。


[Anonymity-0 (Anonymity-0) · GitHub](https://github.com/Anonymity-0)