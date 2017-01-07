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
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	accessFileName := "/var/log/apache2/access_log"
	// Process events
	go func(fileName string) {

		//打开监听的文件
		f, _ := os.Open(fileName)
		defer f.Close()
		rd := bufio.NewReader(f)
		//创建网络链接进行日志传递
		conn, err := net.Dial("tcp", "localhost:9999")
		if err != nil {
			logs.Error("can't connet to localhost:9999")
		}
		wbuf := bufio.NewWriterSize(conn, 1518)
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
						wbuf.WriteString(line)
					}
				} else {
					log.Println("event:", ev.Name)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
				continue
			}
		}
	}(accessFileName)

	err = watcher.Watch(accessFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}

