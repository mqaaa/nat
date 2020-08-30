package main

import (
	"fmt"
	"net"
)

/*
	listener, _ := net.Listen("tcp", "0.0.0.0:9090")
	conn, _ := listener.Accept()
	go func() {
		var buf [128]byte
		n, _ := conn.Read(buf[:])
		defer conn.Close()
	}()

*/

func main() {
	// 1. 监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		fmt.Printf("net Listen, err:%v\n", err)
		return
	}

	// 2. 建立套接字
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("net Accept, err:%v\n", err)
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("net read, err:%v\n", err)
			break
		}
		str := string(buf[:n])
		fmt.Printf("receive from client,date: %v\n", str)
		n, _ = conn.Write([]byte(str))
	}
}
