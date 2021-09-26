/*
@Time : 2021/9/26 16:20
@Author : miwei
@File : main.go
@Description : 描述
*/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	ENV_VERSION = "VERSION"
	HEADER_VERSION = "Version"
)

var (
	flagAddress = flag.String("addr", ":23000", "server listen address")

	// BuildVersion should generate from build script
	BuildVersion = "unknown"
	// BuildDate should generate from build script
	BuildDate = "unknown"
)

func main() {
	fmt.Println("study 20210925", BuildVersion, BuildDate)

	flag.Parse()

	m := &MyIn{}

	go func() {
		m.srv = &http.Server{
			Addr:              *flagAddress,
			Handler:           m,
		}
		fmt.Println("http server listen", *flagAddress)
		err := m.srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("server<listen: %s> running !!!\n", *flagAddress)
	fmt.Println("GET /healthz --> MyIn")
	handleSignal()
	fmt.Println("server stop !!!")

	m.srv.Close()
}

func handleSignal() int {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <- signals
		fmt.Println("get a signal", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// todo: 重新导入配置
			return -1
		default:
			// 退出
			return 0
		}
	}

	return 0
}

type MyIn struct {
	srv *http.Server
}

func (m *MyIn) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 将 request 的 header 写入 response header
	for k, vs := range req.Header {
		fmt.Printf("request header key = %s value = %v\n", k, vs)
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	// 读取环境变量VERSION，写入 response header
	version := os.Getenv(ENV_VERSION)
	if len(version) > 0 {
		w.Header().Set(HEADER_VERSION, version)
	}

	w.WriteHeader(http.StatusOK)

	fmt.Printf("client remote = %s method = %s url = %s code = %d\n", req.RemoteAddr, req.Method, req.URL.Path, http.StatusOK)
}
