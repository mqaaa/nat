package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Recrived request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
	transport := http.DefaultTransport
	// step 1 进行浅靠背，然后新增属性数据
	outReq := new(http.Request)
	*outReq = *req
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		fmt.Printf("ClientIP = %s\n", clientIP)
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step 2 请求下游
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// step 3 把下游数据返回上游
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	_, _ = io.Copy(rw, res.Body)
	_ = res.Body.Close()

}

func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/", &Pxy{})
	_ = http.ListenAndServe("0.0.0.0:8080", nil)
}
