package main

import (
	"fmt"
	"github.com/ShinichR/gotcpServer"

	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Callback struct{}

func (this *Callback) OnConnect() bool {

	fmt.Println("OnConnect:")
	return true
}

type Protocol struct{}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := &gotcpServer.ServerConfig{
		SendChanLimit: 20,
		RecvChanLimit: 20,
		Port:          11125,
	}

	fmt.Println("this is server test example")
	srv := gotcpServer.NewServer(config, &Callback{}, &Protocol{})

	go srv.Start(time.Second)

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()

}
