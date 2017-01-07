package main

import (
	"net"
	"bufio"
	"github.com/astaxie/beego/logs"
	"fmt"
)

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		logs.Error("can't bind port 9999")
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logs.Error(err)
		}
		go handleConnection(conn)
		defer conn.Close()
	}
}

//链接数据处理
//TODO 考虑多个协程输出干扰
func handleConnection(conn net.Conn) {
	rbuf := bufio.NewReader(conn)
	for {
		str, err := rbuf.ReadString('\n')
		if err != nil {
			logs.Error(err)
			break
		}
		fmt.Print(str)
	}
}
