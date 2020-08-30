package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	// 创建连接池
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 链接超时
			KeepAlive: 30 * time.Second, // 长链接超时
		}).DialContext,
		MaxIdleConns:          100, // 最大空闲链接量
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 创建客户端
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	// 访问请求
	resp, err := client.Get("http://127.0.0.1:9090/bye")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bds, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bds))
}
