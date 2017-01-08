package main

import (
	"log"
	"github.com/howeyc/fsnotify"
	"bufio"
	"os"
	"fmt"
	"io"
	"net"
	"github.com/astaxie/beego/logs"
	"compress/zlib"
	"crypto/md5"
	"encoding/hex"
	"flag"
)

var auth_password = ""
var server_host = ""
var accessFileName = "/var/log/apache2/access_log"

func main() {
	p := flag.String("p", "fang", "author password")
	bmd5 := md5.Sum([]byte(*p))
	auth_password = hex.EncodeToString(bmd5[:])
	//init access file
	accessFileName = *flag.String("f", "/var/log/apache2/access_log", "sync file path")
	server_host = *flag.String("s", "localhost", "server host")

	for {
		//loop connect server
		doSync(accessFileName)
		logs.Info("reconnect to server")
	}
}

//进行文件同步操作
func doSync(fileName string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Watch(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	//打开监听的文件
	f, _ := os.Open(fileName)
	//移动指针到末尾
	f.Seek(0, os.SEEK_END)
	defer f.Close()
	rd := bufio.NewReader(f)
	//丢弃一行日志, 防止出现半行情况
	rd.ReadString('\n')
	//创建网络链接进行日志传递
	conn, err := net.Dial("tcp", server_host + ":9999")
	if err != nil {
		logs.Error("can't connet to "+server_host+":9999")
	}

	zwr := zlib.NewWriter(conn)
	//进行认证操作
	n, err := zwr.Write([]byte(auth_password + "\n"));
	if n == 0 || err != nil {
		logs.Error("author error, on write")
		return
	}
	defer zwr.Close()
	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsModify() {
				for {
					line, err := rd.ReadString('\n')
					if err != nil || io.EOF == err {
						break
					}
					fmt.Print(line)
					n, err := zwr.Write([]byte(line))
					if n == 0 || err != nil {
						//TODO 重连
						logs.Error("write data to server error")
						return
					}
					zwr.Flush()
				}
			} else {
				log.Println("event:", ev.Name)
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
			continue
		}
	}
}

