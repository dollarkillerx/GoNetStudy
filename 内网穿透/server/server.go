/**
 * @Author: DollarKillerX
 * @Description: server.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午3:29 2019/12/9
 */
package main

import (
	"log"
	"net"
	"time"
)

func main() {
	listenUdpZero, e := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 9527})
	if e != nil {
		panic(e)
	}
	defer listenUdpZero.Close()
	log.Printf("本机地址: <%s>\n", listenUdpZero.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		i, addr, e := listenUdpZero.ReadFromUDP(data)
		if e != nil {
			log.Printf("error during read: %s", e)
		}
		log.Printf("<%s> %s\n", addr.String(), data[:i])
		peers = append(peers, *addr)
		if len(peers) == 2 {
			log.Printf("进行UDP打洞,建立 %s  < - - > %s 的连接\n", peers[0].String(), peers[1].String())
			listenUdpZero.WriteToUDP([]byte(peers[1].String()), &peers[0])
			listenUdpZero.WriteToUDP([]byte(peers[1].String()), &peers[1])
			time.Sleep(time.Second * 8)
			log.Println("服务器退出,不影响peers间通信")
			return
		}

	}
}
