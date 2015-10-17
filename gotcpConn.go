package gotcpServer

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Error type
var (
	ErrConnClosed    = errors.New("network connection is closed")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface {
	OnConnect() bool

	OnMessage(*Conn, Packet) bool

	OnClose(*Conn)
}

type Conn struct {
	srv           *TcpServer
	conn          *net.TCPConn
	CloseOnce     sync.Once
	CloseFlag     int32
	CloseConnChan chan int
	pktSendChan   chan Packet
	pktRecvChan   chan Packet
	connwaitGroup *sync.WaitGroup
}

// newConn returns a wrapper of raw conn
func newConn(conn *net.TCPConn, srv *TcpServer) *Conn {
	return &Conn{
		srv:           srv,
		conn:          conn,
		CloseConnChan: make(chan int),
		pktSendChan:   make(chan Packet, srv.config.SendChanLimit),
		pktRecvChan:   make(chan Packet, srv.config.SendChanLimit),
		connwaitGroup: &sync.WaitGroup{},
	}
}

// gotcpServer  conn do function

func (c *Conn) Do() {
	defer c.srv.waitGroup.Done()
	c.connwaitGroup.Add(3)
	if !c.srv.callback.OnConnect() {
		return
	}
	go c.pktReadDo()
	go c.pktHandDo()
	go c.pktWriteDo()
	c.connwaitGroup.Wait()
	fmt.Println("conn is closed...")
}

func (c *Conn) isClosed() bool {
	return atomic.LoadInt32(&c.CloseFlag) == 1
}

func (c *Conn) Close() {

	c.CloseOnce.Do(func() {
		atomic.StoreInt32(&c.CloseFlag, 1)
		close(c.CloseConnChan)
		close(c.pktRecvChan)
		close(c.pktSendChan)
		c.conn.Close()
		c.srv.callback.OnClose(c)
	})
}

func (c *Conn) WritePkt(p Packet, timeout time.Duration) (err error) {
	if c.isClosed() {
		return ErrConnClosed
	}

	defer func() {
		if e := recover(); e != nil {
			err = ErrConnClosed
		}
	}()

	if timeout == 0 {
		select {
		case c.pktSendChan <- p:
			return nil
		default:
			return ErrWriteBlocking
		}
	} else {
		select {
		case c.pktSendChan <- p:
			return nil

		case <-c.CloseConnChan:
			return ErrConnClosed
		case <-time.After(timeout):
			return ErrWriteBlocking
		}
	}

}

func (c *Conn) pktReadDo() {
	defer func() {
		recover()
		c.Close()
		c.connwaitGroup.Done()
		fmt.Println("pktReadDo done")
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return
		case <-c.CloseConnChan:
			return

		default:
		}

		p, err := c.srv.protocol.ReadPacket(c.conn)
		if err != nil {
			return
		}
		c.pktRecvChan <- p

	}

}

func (c *Conn) pktWriteDo() {
	defer func() {
		recover()
		c.Close()
		c.connwaitGroup.Done()
		fmt.Println("pktWriteDo done")
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return
		case <-c.CloseConnChan:
			return
		case p := <-c.pktSendChan:
			if c.isClosed() {
				return
			}
			if _, err := c.conn.Write(p.Serialize()); err != nil {
				return
			}
		default:
		}

	}

}

func (c *Conn) pktHandDo() {
	defer func() {
		recover()
		c.Close()
		c.connwaitGroup.Done()
		fmt.Println("pktHandDo done")
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.CloseConnChan:
			return

		case p := <-c.pktRecvChan:
			if c.isClosed() {
				return
			}
			if !c.srv.callback.OnMessage(c, p) {
				return
			}
		}
	}
}
