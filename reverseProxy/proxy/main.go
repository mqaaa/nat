package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	proxyAddr = "http://127.0.0.1:2003"
	port = "2002"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port " + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(writer http.ResponseWriter, request *http.Request) {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		log.Println(err)
		return
	}
	request.URL.Scheme = proxy.Scheme
	request.URL.Host = proxy.Host

	// 请求下游
	transport := http.DefaultTransport
	res, err := transport.RoundTrip(request)
	if err != nil {
		log.Println(err)
		return
	}

	// step 3 把下游数据返回上游
	for key, value := range res.Header {
		for _, v := range value {
			writer.Header().Add(key, v)
		}
	}
	defer res.Body.Close()
	//bufio.NewReader(res.Body).WriteTo(writer)
	_, _ = io.Copy(writer, res.Body)
}
