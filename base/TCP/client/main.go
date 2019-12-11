/**
 * @Author: DollarKillerX
 * @Description: tcp连接客户端
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 上午11:49 2019/12/10
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr, e := net.ResolveTCPAddr("tcp", ":9999")
	if e != nil {
		panic(e)
	}
	// 拨号 本机,远程的机器
	conn, e := net.DialTCP("tcp", nil, addr)
	if e != nil {
		panic(e)
	}
	defer conn.Close()
	fmt.Println(conn.LocalAddr().String() + " : Client connected!")

	go writeMsg(conn)
	readMsg(conn)
}

// 读协程
func readMsg(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("关闭连接")
				conn.Close()
			} else {
				log.Println(err)
				continue
			}
		}
		log.Println()
		log.Println(string(buf[:n]))
	}
}

func writeMsg(conn *net.TCPConn) {
	for {
		fmt.Printf("请输入msg:")
		msg := ""
		_, err := fmt.Scanf("%s", &msg)
		if err != nil {
			panic(err)
		}

		conn.Write([]byte(msg))
	}
}
