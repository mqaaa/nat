package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RuleServer struct {
	Addr string
}

func main() {
	rs1 := &RuleServer{Addr: "127.0.0.1:2003"}
	rs1.run()
	rs2 := &RuleServer{Addr: "127.0.0.1:2004"}
	rs2.run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (r *RuleServer) run() {
	log.Println("Starting httpServer as " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	Server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}
	go func() {
		log.Fatal(Server.ListenAndServe())
	}()
}

func (r *RuleServer) HelloHandler(writer http.ResponseWriter, request *http.Request) {
	Path := fmt.Sprintf("http://%s%s\n", r.Addr, request.URL.Path)
	_, _ = io.WriteString(writer, Path)
}

func (r *RuleServer) ErrorHandler(writer http.ResponseWriter, request *http.Request) {
	Path := "error header"
	writer.WriteHeader(500)
	_, _ = io.WriteString(writer, Path)
}
