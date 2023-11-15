
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
#### inet_addr

```c
#include <arpa/inet.h> 
in_addr_t inet_addr(const char * string); 
	//成功时返回32位大端序整数型值,失败时返回 INADDR_NONE。

```
如果向该函数传递类似“211.214.107.99”的点分十进制格式的字符串,它会将其转换为32 位整数型数据并返回。当然,该整数型值满足网络字节序。另外,该函数的返回值类型in_addr_t 在内部声明为32位整数型。


在 Go 中，你可以使用 `net` 包的 `ParseIP` 函数来解析 IP 地址。这个函数可以处理 IPv4 和 IPv6 地址，如果输入的字符串不是有效的 IP 地址，它会返回 `nil`。然后，你可以使用 `encoding/binary` 包的 `BigEndian.Uint32` 或 `LittleEndian.Uint32` 函数将 `net.IP` 类型的 IP 地址转换为网络字节序或主机字节序的整数。

```go
package main

import (
    "encoding/binary"
    "fmt"
    "net"
)

func main() {
    ip := net.ParseIP("1.2.3.4")
    if ip == nil {
        fmt.Println("Invalid IP address")
        return
    }
    ip = ip.To4()
    if ip == nil {
        fmt.Println("Not an IPv4 address")
        return
    }
    fmt.Printf("IP as integer (network byte order): %x\n", binary.BigEndian.Uint32(ip))
    fmt.Printf("IP as integer (host byte order): %x\n", binary.LittleEndian.Uint32(ip))
}
```

inter_addr.go

```go
package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const INADDR_NONE = 0xffffffff

// inet_addr converts an IPv4 address in dot-decimal notation into a 32-bit integer in network byte order.
// If the IP address is invalid, it returns INADDR_NONE (0xffffffff).
func inet_addr(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0xffffffff
	}
	ip = ip.To4()
	if ip == nil {
		return 0xffffffff
	}
	return binary.BigEndian.Uint32(ip)

}
func main() {
	addr1 := "1.2.3.4"
	addr2 := "1.2.3.256"

	conv_addr := inet_addr(addr1)

	if conv_addr == INADDR_NONE {
		fmt.Println("Error occured!")
	} else {
		fmt.Printf("Network ordered integer addr: %#x\n", conv_addr)
	}
	conv_addr = inet_addr(addr2)
	if conv_addr == INADDR_NONE {
		fmt.Println("Error occured!")
	} else {
		fmt.Printf("Network ordered integer addr: %#x\n", conv_addr)
	}
}



```

#### inet_aton

`inet_aton`函数 与inet_addr函数在功能上完全相同,也将字符串形式IP地址转换为32位网络字节序整数并返回。只不过该函数利用了in_addr结构体,且其使用频率更高。


```c
#include <arpa/inet.h> 
int inet_aton(const char * string,struct in_addn * addn); 
//成功时返回1(true),失败时返回 0(false)。
//string 含有需转换的IP地址信息的字符串地址值。
//addr 将保存转换结果的in_addr结构体变量的地址值。

```

inet_aton.go

```go
package main

import (
    "encoding/binary"
    "fmt"
    "log"
    "net"
)

// inet_aton converts an IPv4 address in dot-decimal notation into a 32-bit integer in network byte order.
// If the IP address is valid, it stores the result in the given *uint32 and returns true.
// If the IP address is invalid, it returns false.
func inet_aton(ipStr string, ip *uint32) bool {
    parsedIP := net.ParseIP(ipStr)
    if parsedIP == nil {
        return false
    }
    parsedIP = parsedIP.To4()
    if parsedIP == nil {
        return false
    }
    *ip = binary.BigEndian.Uint32(parsedIP)
    return true
}

func main() {
    addr := "127.232.124.79"
    var ip uint32

    if !inet_aton(addr, &ip) {
        log.Fatalln("Conversion error")
    } else {
        fmt.Printf("Network ordered integer addr: %#x\n", ip)
    }
}

```

#### inet_aton
上述运行结果无关紧要,更重要的是大家要熟练掌握该函数的调用方法。最后再介绍一个与inet_aton函数正好相反的函数,此函数可以把网络字节序整数型IP地址转换成我们熟悉的字符串形式。

```c
#include <arpa/inet.h> 
char * inet_ntoa(struct in_addr adr);
//成功时返回转换的字符串地址值,失败时返回-1。

```
该函数将通过参数传入的整数型IP地址转换为字符串格式并返回。
但调用时需小心,返回值类型为char指针。返回字符串地址意味着字符串已保存到内存空间,但该函数未向程序员要求分配内存,而是在内部申请了内存并保存了字符串。也就是说,调用完该函数后,应立即将字符串信息复制到其他内存空间。总之,再次调用inet_ntoa函数前返回的字符串地址值是有效的。若需要长期保存,则应将字符串复制到其他内存空间。

Go语言标准库中并没有直接提供对C语言中的网络地址转换函数的封装。如果要在Go中实现类似的功能，可以使用net包中的IP和IPv4类来进行IP地址的转换和操作。以下是一个简单的示例代码：

```go
package main

import (
	"fmt"
	"net"
)

func main1() {
	var addr1, addr2 uint32 = 0x1020304, 0x1010101

	ip1 := net.IPv4(byte(addr1>>24), byte(addr1>>16), byte(addr1>>8), byte(addr1))
	ip2 := net.IPv4(byte(addr2>>24), byte(addr2>>16), byte(addr2>>8), byte(addr2))

	fmt.Printf("Dotted-Decimal notation1: %s \n", ip1.String())
	fmt.Printf("Dotted-Decimal notation2: %s \n", ip2.String())
}


```

用go自行写的 inet_aton函数

```go
package main

import (
	"fmt"
	"net"
)

// inet_ntoa converts a 32-bit integer in network byte order into a dotted-decimal IP address.
func inet_ntoa(ipInt uint32) string {
	ipBytes := make([]byte, 4)
	ipBytes[0] = byte(ipInt >> 24)
	ipBytes[1] = byte(ipInt >> 16)
	ipBytes[2] = byte(ipInt >> 8)
	ipBytes[3] = byte(ipInt)
	return net.IP(ipBytes).String()
}

func main() {
	var addr1, addr2 uint32 = 0x1020304, 0x1010101

	fmt.Printf("Dotted-Decimal notation1: %s \n", inet_ntoa(addr1))
	fmt.Printf("Dotted-Decimal notation2: %s \n", inet_ntoa(addr2))
}


```
### 网络地址初始化
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.5lfoz23nb62o.webp)

上述代码中,memset函数将每个字节初始化为同一值:第一个参数为结构体变量addr的地址值,即初始化对象为addr;第二个参数为0,因此初始化为0;最后一个参数中传入addr的长度,因此addr的所有字节均初始化为0。这么做是为了将sockaddr_in结构体的成员sin_zero初始化为0。
另外,最后一行代码调用的atoi函数把字符串类型的值转换成整数型。总之,上述代码利用字符串格式的IP地址和端口号初始化了sockaddr_in结构体变量。
另外,代码中对IP地址和端口号进行了硬编码,这并非良策,因为运行环境改变就得更改代码。因此,我们运行示例main函数时传入IP地址和端口号。
### 客户端地址信息初始化

上述网络地址信息初始化过程主要针对服务器端而非客户端。给套接字分配IP地址和端口号主要是为下面这件事做准备:

	“请把进入IP 211.217.168.13、9190端口的数据传给我!”

反观客户端中连接请求如下: 

	“请连接到IP 211.217.168.13、9190端口!”

请求方法不同意味着调用的函数也不同。服务器端的准备工作通过`bind`函数完成,而客户端则通过`connect`函数完成。因此,函数调用前需准备的地址值类型也不同。服务器端声明`sockaddr_in` 结构体变量,将其初始化为赋子**服务器端IP和套接字的端口号**,然后调用**bind函数**;而客户端则声明sockaddr_in结构体,并初始化为要与之连接的服务器端套接字的IP和端口号,然后调用**connect函数**。

### INADDR_ANY
![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.4hipk7niksg0.webp)
与之前方式最大的区别在于,利用常数**INADDR_ANY**分配服务器端的IP地址。若采用这种方式,则可**自动获取运行服务器端的计算机IP地址**,不必亲自输人。而且,若同一计算机中已分配多个IP地址(多宿主(Multi-homed)计算机,一般路由器属于这一类),则只要端口号一致, 就可以从不同IP地址接收数据。因此,服务器端中优先考虑这种方式。而客户端中除非带有一部分服务器端功能,否则不会采用。

初始化服务器端套接字时应分配所属计算机的IP地址,因为初始化时使用的IP地址非常明确,那为何还要进行IP初始化呢?如前所述,同一计算机中可以分配多个IP地址, **实际IP地址的个数与计算机中安装的NIC的数量相等**。即使是**服务器端套接字,也需要决定应接收哪个IP传来的(哪个NIC传来的)数据**。因此,服务器端套接字初始化过程中要求IP地址信息。另外,若只有1个NIC,则直接使用INADDR_ANY。

在 Go 中，你可以使用空字符串 `""` 作为 IP 地址来代表 INADDR_ANY，这表示监听所有的 IP 地址。以下是一个简单的 TCP 服务器示例，它监听所有的 IP 地址和一个特定的端口：

```go
package main

import (
    "log"
    "net"
)

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    // Handle the connection
    defer conn.Close()
}

```
在这个示例中，`net.Listen("tcp", ":8080")` 会监听所有的 IP 地址和端口 8080。当有新的连接到来时，`listener.Accept()` 会返回一个新的 `net.Conn` 对象，然后你可以在新的 goroutine 中处理这个连接。


### 向套接字分配网络地址

既然已讨论了sockaddr_in结构体的初始化方法,接下来就把初始化的地址信息分配给套接字。bind函数负责这项操作。

```c
#include <sys/socket.h> 
int bind(int sockfd, struct sockaddr * myaddr, socklen_t addrLen);
//成功时返回0,失败时返回-1。
//sockfd 要分配地址信息(IP地址和端口号)的套接字文件描述符。
//myaddr 存有地址信息的结构体变量地址值。
//addrlen 第二个结构体变量的长度。
```

在 Go 语言中，`Bind` 函数是 `syscall` 包中的一个函数，用于将本地协议地址 `addr` 绑定到文件描述符 `fd`。函数原型如下：

```go
func Bind(fd int, sa Sockaddr) (err error) {
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return bind(fd, ptr, n)
}
```

其中，`fd` 是通过 `Socket` 函数获取的文件描述符，`addr` 是一个实现了 `Sockaddr` 接口的网络地址。

两者的主要区别在于：

1. Go 的 `Bind` 函数使用了接口 `Sockaddr`，这使得你可以传入任何实现了 `Sockaddr` 接口的类型，如 `SockaddrInet4`、`SockaddrInet6`、`SockaddrUnix` 等。而 C 的 `bind` 函数需要一个指向 `struct sockaddr` 的指针，需要手动进行类型转换。

2. Go 的 `Bind` 函数返回一个错误值，你可以直接检查这个错误值来确定 `Bind` 函数是否成功。而 C 的 `bind` 函数返回一个整数，需要检查这个整数和 `errno` 来确定 `bind` 函数是否成功。

3. Go 的 `Bind` 函数处理了一些底层的细节，如网络字节序的转换。而在 C 中，需要手动进行这些操作。


```go
func bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error) {
	_, _, e1 := syscall(abi.FuncPCABI0(libc_bind_trampoline), uintptr(s), uintptr(addr), uintptr(addrlen))
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
```

在 Go 语言的 `syscall` 包中，`Bind` 和 `bind` 函数都是用来将本地协议地址绑定到文件描述符的。
`bind` 函数是一个私有的函数，它直接接受一个文件描述符和一个 `unsafe.Pointer` 类型的地址作为参数，然后调用系统调用 `bind`。这个函数通常不会直接被用户代码调用，而是被 `Bind` 函数调用。
`bind` 函数是通过 `syscall` 函数调用 `libc_bind_trampoline` 函数来实现的，这个函数是一个跳板函数，它会跳转到动态链接库中的 `bind` 函数。
总的来说，`Bind` 函数提供了一个更高级的接口，它处理了类型转换和错误处理，而 `bind` 函数是一个更底层的接口，它直接调用系统调用。


## 3.5 基于Windows的实现
略

## 3.6 习题
1. IP地址族IPv4和IPv6有何区别?在何种背景下诞生了IPv6? 

	IPv4与IPv6的差别主要是表示IP地址所用的字节数,目前通用的地址族为IPv4。IPv6是为了应对2010年前后IP地址耗尽的问题而提出的标准。

2. 通过IPv4网络ID、主机ID及路由器的关系说明向公司局域网中的计算机传输数据的过程。

	网络地址(网络ID)是为区分网络而设置的一部分IP地址。假设向WWW.SEMI.COM公司传输数据,该公司内部构建了局域网,把所有计算机连接起来。因此,首先应向SEMI.COM网络传输数据,也就是说,并非一开始就浏览所有4字节IP地址,进而找到目标主机;而是仅浏览4字节IP地址的网络地址,先把数据传到SEMI.COM的网络。SEMI.COM网络(构成网络的路由器)接收到数据后,浏览传输数据的主机地址(主机ID)并将数据传给目标计算机。

3. 套接字地址分为IP地址和端口号。为什么需要地址和端口号?或者说,通过IP可以区分哪些对象?通过端口号可以区分哪些对象?

	套接字地址分为IP地址和端口号，是为了在网络中唯一标识一个通信端点。**IP地址**用于区分不同的主机，即不同的计算机。**端口号**用于区分同一主机上的不同进程，即不同的应用程序。

4. 请说明IP地址的分类方法,并据此说出下面这些IP地址的分类。
	![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.19gzb4mtinb4.png)
	- 214.121.212.102 （C类）
	- 120.101.122.89 （A类）
	- 129.78.102.211 （B类）

5. 计算机通过路由器或交换机连接到互联网。请说出路由器和交换机的作用。

	若想构建网络,需要一种物理设备完成外网与本网主机之间的数据交换,这种设备便是路由器或交换机。
	
6. 什么是知名端口?其范围是多少?知名端口中具有代表性的HTTP和FTP端口号各是多少? 
	
	0-1023是知名端口(Well-known PORT),一般分配给特定应用程序。HTTP的端口号是**80**，FTP的端口号是**21**

7. 向套接字分配地址的bind函数原型如下: `int bind(int sockfd, struct sockaddr *myaddr,socklen_t addrlen);` 而调用时则用`bind(serv_sock,(struct sockaddr *)&sery_addr, sizeof (serv_addr));` 此处`serv_addr`为`sockaddr_in`结构体变量。与函数原型不同,传入的是`sockaddr_in`结构体变量,请说明原因。

	`sockaddr_in` 结构体是 `sockaddr` 结构体的一种特定类型，因此可以通过类型转换将其传递给 `bind` 函数。这是因为 `sockaddr_in` 结构体包含了 `sockaddr` 结构体的所有成员，所以在实际调用中可以将 `sockaddr_in` 结构体的指针转换为 `sockaddr` 结构体的指针，从而符合 `bind` 函数的参数要求。

8. 请解释大端序、小端序、网络字节序,并说明为何需要网络字节序。

	- **大端序（Big-Endian）**：数据的低位字节存储在内存的低地址，高位字节存储在内存的高地址。
		![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.27xmkxdgjg74.webp)
	- **小端序（Little-Endian）**：数据的高位字节存储在内存的低地址，低位字节存储在内存的高地址。
		![image](https://cdn.statically.io/gh/Anonymity-0/Picgo@note_picture/img/image.2gkyxou8dm4g.webp)
	- **网络字节序（Network Byte Order）**：是指在网络传输中采用的字节序。网络字节序采用大端序，这是因为大端序与人类阅读数字的顺序一致，因此更容易理解和识别。
	- 网络字节序是网络传输的标准，因此在网络传输中采用网络字节序可以**确保数据在不同计算机之间正确传输**。

9. 大端序计算机希望把4字节整数型数据12传递到小端序计算机。请说出数据传输过程中发生的字节序变换过程。

	数据 12 在网络传输过程中的字节序没有发生变化，仍然是大端序。小端序计算机在接收到数据后，需要将数据中的高位字节和低位字节进行交换，以将数据转换为小端序。
	
| 大端序       | 网络字节序 | 小端序     |
| ---------- | ---------- | ---------- |
|    0x0000000C     |      0x0000000C      |    0x0C000000         |

10. 怎样表示回送地址?其含义是什么?如果向回送地址传输数据将发生什么情况?
	- 回送地址（loopback address）是指**本地主机的 IP 地址**。在 IPv4 中，回送地址为 **127.0.0.1**。
	- 回送地址用于本地主机之间的通信。如果向回送地址传输数据，则数据将会被本地主机接收并处理。



