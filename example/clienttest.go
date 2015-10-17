package main

import (
	"bufio"
	"fmt"
	"github.com/ShinichR/gotcpServer/protocol"
	"log"
	"net"
	"os"
)

var ReadStdinChan chan string
var ReadServerChan chan string

func ReadFromStdin() {
	for {
		fmt.Print("input:")
		cmdReader := bufio.NewReader(os.Stdin)
		cmdStr, err := cmdReader.ReadString('\n')
		if err == nil {
			ReadStdinChan <- cmdStr[:len(cmdStr)-1]
		}

	}
}

func ReadFromServer(c *net.TCPConn) {
	for {
		echoProtocol := &protocol.EchoProtocol{}

		p, err := echoProtocol.ReadPacket(c)
		if err == nil {
			ReadServerChan <- string(p.(*protocol.EchoPacket).GetBody())
		}
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:11125")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	fmt.Println("connected server....")

	ReadStdinChan = make(chan string)
	ReadServerChan = make(chan string)
	go ReadFromStdin()
	go ReadFromServer(conn)

	for {
		select {
		case p := <-ReadStdinChan:
			fmt.Println("read from stdin:", p)
			//

			conn.Write(protocol.NewEchoPacket([]byte(p), false).Serialize())

		case v := <-ReadServerChan:
			fmt.Printf("Server reply:\n", v)

		default:

		}

	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
