/**
 * @Author: DollarKillerX
 * @Description: main.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午3:59 2019/12/9
 */
package main

import (
	"github.com/dollarkillerx/obfuscation/obbase64"
	"github.com/dollarkillerx/obfuscation/zipstr"
	"log"
	"net"
)

var key = "45,83,21,69,70,79,16,87,90,04,99,65,04,10,06,09"

func main() {
	addr, e := net.ResolveUDPAddr("udp", "0.0.0.0:8081")
	if e != nil {
		panic(e)
	}
	conn, e := net.ListenUDP("udp", addr)
	if e != nil {
		panic(e)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		i, udpAddr, e := conn.ReadFromUDP(buf)
		if e != nil {
			log.Println("err: ")
			log.Println(e)
			continue
		}
		msg := string(buf[:i])
		unzip := zipstr.Unzip(msg)
		log.Println(msg)
		bytes, e := obbase64.Decode(key, unzip)
		if e != nil {
			log.Println("err: ")
			log.Println(e)
			continue
		}
		log.Printf("Ta <%s> 发送来的Msg: %s\n", udpAddr.String(), string(bytes))

		result := "已经收到来自你的数据: " + udpAddr.String()
		s, e := obbase64.Encode(key, []byte(result))
		if e != nil {
			log.Println("err: ")
			log.Println(e)
			continue
		}
		zip := zipstr.Zip(s)
		_, e = conn.WriteToUDP([]byte(zip), udpAddr)
		if e != nil {
			log.Println("err: ")
			log.Println(e)
			continue
		}
	}
}
