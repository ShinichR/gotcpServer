package gotcpServer

import (
	"net"
)

type Packet interface {
	Serialize() []byte
}

type LayerProtocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}
