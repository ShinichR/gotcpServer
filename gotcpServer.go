package gotcpServer


import (
	"net"
	"sync"
	"time"
	"log"
	"fmt"
)


type ConnCallback struct{

}

type LayerProtocol struct{

}

type ServerConfig struct{
	SendChanLimit    uint32
	RecvChanLimit    uint32
	port             uint32
}

type TcpServer struct{
	config *ServerConfig
	callback ConnCallback
	protocol LayerProtocol
	exitChan chan int
	waitGroup *sync.WaitGroup
	
}

func NewServer(config *ServerConfig , callback ConnCallback , protocol LayerProtocol) *TcpServer{
	return &TcpServer{
		config: config,
		callback: callback,
		protocol:protocol,
		exitChan: make(chan int),
		waitGroup: &sync.WaitGroup{},
	}

}

func (s *TcpServer) Start(acTimeout time.Duration){ 
	tcpAddr, err := net.ResolveTCPAddr("tcp4",":"+string(s.config.port))
	Errdeal(err)
	listener,err := net.ListenTCP("tcp",tcpAddr)
	Errdeal(err)
	
	s.waitGroup.Add(1)
	defer func(){
		listener.Close()
		s.waitGroup.Done()
	}()
	
	for{
		
		select{
			case <-s.exitChan:
				return
			default:
		}
		listener.SetDeadline(time.Now().Add(acTimeout))
		_,err := listener.AcceptTCP()
		if err != nil{
			continue
		}
		
		s.waitGroup.Add(1)
	        fmt.Println("new client connected...")	

	}
}
func (s *TcpServer) Stop(){
	close(s.exitChan)
	s.waitGroup.Wait()
}
func Errdeal(err error){
	if err != nil{
	      log.Fatal("err",err)
	}
}


