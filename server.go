package main

import (
	"net"
	"bufio"
	"github.com/astaxie/beego/logs"
	"fmt"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"flag"
)

var auth_password = ""

func main() {
	p := flag.String("p", "fang", "author password")
	bmd5 := md5.Sum([]byte(*p))
	auth_password = hex.EncodeToString(bmd5[:])

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
	r, _ := zlib.NewReader(conn)
	rbuf := bufio.NewReader(r)

	str, _ := rbuf.ReadString('\n')
	if !author(str) {
		logs.Error("token error, conn close")
		conn.Close()
		return
	}
	logs.Info("autho success")

	for {
		str, err := rbuf.ReadString('\n')
		if err != nil {
			logs.Error(err)
			break
		}

		fmt.Print(str)
	}
}

func author(token string) (bool) {
	return auth_password + "\n" == token
}
