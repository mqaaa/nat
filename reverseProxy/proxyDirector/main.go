package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var Addr = "127.0.0.1:2002"

func main() {
	rs1 := "http://127.0.0.1:2003/base"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		panic(err1)
	}
	proxy := NewSingleHostReverseProxy(url1)
	log.Fatal(http.ListenAndServe(Addr, proxy))
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = SingleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
		req.Header.Set("X-Real-Ip", "11.11.11.11")
	}
	modifyFunc := func(res *http.Response) error {
		if res.StatusCode != 200 {
			oldPayLoad, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			newPayLoad := []byte("hello " + string(oldPayLoad))
			res.Body = ioutil.NopCloser(bytes.NewBuffer(newPayLoad))
			res.ContentLength = int64(len(newPayLoad))
			res.Header.Set("Content-Length", fmt.Sprint(len(newPayLoad)))
		}
		return nil
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc}
}

func SingleJoiningSlash(path, path2 string) string {
	aslash := strings.HasSuffix(path, "/")
	bslash := strings.HasPrefix(path2, "/")
	switch {
	case aslash && bslash:
		return path + path2[1:]
	case !aslash && !bslash:
		return path + "/" + path2
	}
	return path + path2
}
