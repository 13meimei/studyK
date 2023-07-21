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
	"net/http/pprof"
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

	m := NewMyIn()
	m.SetMux()
	srv := &http.Server{
		Addr:              *flagAddress,
		Handler:           m.GetMux(),
	}
	m.SetHttpServer(srv)

	go func() {
		fmt.Println("http server listen", *flagAddress)
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("server<listen: %s> running !!!\n", *flagAddress)
	fmt.Println("GET /healthz --> MyIn.Healthz")
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
	mux *http.ServeMux
}

func NewMyIn() *MyIn {
	return &MyIn{}
}

func (m *MyIn) SetHttpServer(srv *http.Server) {
	m.srv = srv
}

func (m *MyIn) SetMux() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	mux.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	mux.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	mux.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	mux.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)

	mux.HandleFunc("/healthz", m.Healthz)

	m.mux = mux
}

func (m *MyIn) GetMux() *http.ServeMux {
	return m.mux
}

func (m *MyIn) Healthz(w http.ResponseWriter, req *http.Request) {
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
