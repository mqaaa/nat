package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

/*
	conn, err := net.Dial("tcp", "localhost:9090")
	for {
		_, err = conn.Write([]byte(trimInput))
	}
	_ = conn.Close()
*/

func main() {
	// 1. 链接服务器
	conn, err := net.Dial("tcp", "localhost:9090")
	defer func() {
		_ = conn.Close()
		time.Sleep(30 * time.Second)
	}()
	if err != nil {
		fmt.Printf("connect failed, err :%v\n", err.Error())
		return
	}
	go listen(conn)
	// 2. 读取命令行输入
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err :%v\n", err.Error())
			break
		}
		trimInput := strings.TrimSpace(input)
		if trimInput == "Q" {
			break
		}
		_, err = conn.Write([]byte(trimInput))
		if err != nil {
			fmt.Printf("Write failed, err :%v\n", err.Error())
			break
		}
	}

}

// 小课堂哈哈哈哈
func listen(conn net.Conn) {
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("net read, err:%v\n", err)
		}
		str := string(buf[:n])
		fmt.Printf("receive from client,date: %v\n", str)
	}
}
