/**
 * @Author: DollarKillerX
 * @Description: main.go
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午4:07 2019/12/9
 */
package main

import (
	"fmt"
	"github.com/dollarkillerx/obfuscation/obbase64"
	"github.com/dollarkillerx/obfuscation/zipstr"
	"log"
	"net"
	"sync"
)

func main() {
	conn, e := net.Dial("udp", "0.0.0.0:8081")
	if e != nil {
		panic(e)
	}
	defer conn.Close()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go read(conn, &wg)
	go write(conn, &wg)

	wg.Wait()
}

var key = "45,83,21,69,70,79,16,87,90,04,99,65,04,10,06,09"

func write(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg := ""
		fmt.Scanf("%s", &msg)
		var err error
		msg, err = obbase64.Encode(key, []byte(msg))
		if err != nil {
			log.Println("err: ")
			log.Println(err)
			continue
		}
		zip := zipstr.Zip(msg)
		conn.Write([]byte(zip))
	}
}

func read(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, 1024)
	for {
		n, e := conn.Read(buf)
		if e != nil {
			log.Println("err: ")
			log.Println(e)
			break
		}

		bytes, e := obbase64.Decode(key, zipstr.Unzip(string(buf[:n])))
		if e != nil {
			log.Println("err: ")
			log.Println(e)
			continue
		}
		log.Println("data: ", string(bytes))
	}
}
