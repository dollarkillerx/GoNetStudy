/**
 * @Author: DollarKillerX
 * @Description: main.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 上午11:34 2019/12/10
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var close = "Close 1013 8751"

func main() {
	addr, e := net.ResolveTCPAddr("tcp", ":9999")
	if e != nil {
		panic(e)
	}

	listener, e := net.ListenTCP("tcp", addr)
	if e != nil {
		panic(e)
	}
	defer listener.Close()
	for {
		conn, e := listener.AcceptTCP() // 接受一个新的连接
		if e != nil {
			log.Println(e)
			continue
		}
		fmt.Println("A client connected :" + conn.RemoteAddr().String())
		go tcpPipe(conn)
	}
}

// 处理每一个用户请求
func tcpPipe(conn *net.TCPConn) {
	// tcp连接地址
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("Disconnected:" + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn) // 建立读取器
	buf := make([]byte, 1024)       // 缓冲区
	for {
		n, err := reader.Read(buf)
		if err != nil || err == io.EOF {
			log.Println(err)
			break
		}

		data := string(buf[:n])
		if data == close {
			log.Printf("客户端%s 请求关闭连接\n", ipStr)
			return
		}
		fmt.Println("client msg: ", data)

		msg := time.Now().String() + ipStr + " Server Say Hello! \n"
		conn.Write([]byte(msg))
	}
}
