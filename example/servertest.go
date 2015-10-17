package main

import (
	"fmt"
	"github.com/ShinichR/gotcpServer"
	"github.com/ShinichR/gotcpServer/protocol"
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

func (this *Callback) OnMessage(c *gotcpServer.Conn, p gotcpServer.Packet) bool {
	echoPacket := p.(*protocol.EchoPacket)
	fmt.Printf("OnMessage:%v\n", string(echoPacket.GetBody()))
	c.WritePkt(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)

	return true
}

func (this *Callback) OnClose(c *gotcpServer.Conn) {

	fmt.Println("OnClose:")

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := &gotcpServer.ServerConfig{
		SendChanLimit: 20,
		RecvChanLimit: 20,
		Port:          11125,
	}

	fmt.Println("this is server test example")
	srv := gotcpServer.NewServer(config, &Callback{}, &protocol.EchoProtocol{})

	go srv.Start(time.Second)

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()

}
