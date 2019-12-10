/**
 * @Author: DollarKillerX
 * @Description: 内网穿透客户端
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 上午10:27 2019/12/10
 */
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const HAND_SHAKE_MSG = "我是打洞消息"

func main() {
	if len(os.Args) == 0 {
		fmt.Println("请输入 一个客户端表示")
		os.Exit(0)
	}
	tag := os.Args[1]
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9901}
	dstAddr := &net.UDPAddr{IP: net.ParseIP(tag), Port: 9527}
	conn, e := net.DialUDP("udp", srcAddr, dstAddr)
	if e != nil {
		log.Panic(e)
	}
	defer conn.Close()
	if _, e := conn.Write([]byte("Hello,I'm new peer:" + tag)); e != nil {
		log.Panic(e)
	}
	buf := make([]byte, 1024)
	i, addr, e := conn.ReadFromUDP(buf)
	if e != nil {
		log.Panic(e)
	}

	// 获取建立连接的内网机器
	udpAddr := parserAddr(string(buf[:i]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, addr, udpAddr.String())
	// 开始打洞
	bidirectionHole(srcAddr, &udpAddr)
}

// 打洞 本机和建立连接的内网机器
func bidirectionHole(srcAddr *net.UDPAddr, antherAddr *net.UDPAddr) {
	conn, e := net.DialUDP("udp", srcAddr, antherAddr)
	if e != nil {
		log.Panic(e)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, e := conn.Write([]byte(HAND_SHAKE_MSG)); e != nil {
		log.Println("Send handshake:", e)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err := conn.Write([]byte("from[" + srcAddr.String())); err != nil {
				log.Println("seed msg fail", err)
			}
		}
	}()

	for {
		buf := make([]byte, 1024)
		i, _, e := conn.ReadFromUDP(buf)
		if e != nil {
			log.Println("error during read: %s\n", e)
		} else {
			log.Printf("接受到的数据是:%s\n", buf[:i])
		}
	}

}

// 将传入的udp str 转为 udp结构体
func parserAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	i, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: i,
	}
}
