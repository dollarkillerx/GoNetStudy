/**
 * @Author: DollarKillerX
 * @Description: golang 实现http or https代理服务
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 上午11:31 2019/12/10
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/dollarkillerx/easyutils/clog"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {
	// 监听服务
	listener, e := net.Listen("tcp", ":8088")
	if e != nil {
		log.Panic(e)
	}
	for {
		// 接受每个用户的请求
		conn, e := listener.Accept()
		if e != nil {
			log.Println(e)
			continue
		}
		go handleClientRequest(conn)
	}
}

func handleClientRequest(conn net.Conn) {
	if conn == nil {
		return
	}
	defer func() {
		conn.Close()
		log.Println("连接被释放了")
	}()

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		clog.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortUrl, err := url.Parse(strings.TrimSpace(host))
	if err != nil {
		clog.Println(err)
		return
	}

	if hostPortUrl.Opaque == "443" { // https访问
		address = hostPortUrl.Scheme + ":443"
	} else { // http
		if strings.Index(hostPortUrl.Host, ":") == -1 {
			// host不带端口默认 80
			address = hostPortUrl.Host + ":80"
		} else {
			address = hostPortUrl.Host
		}
	}

	// 与目标通讯地址建立连接
	server, err := net.Dial("tcp", address)
	if err != nil {
		clog.Println(err)
		return
	}

	// 数据转发
	if method == "CONNECT" {
		fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(b[:n])
	}

	// 进行转发
	go io.Copy(server, conn)
	io.Copy(conn, server)
}