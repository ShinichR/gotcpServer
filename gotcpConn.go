package gotcpServer

import ()

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface {
	OnConnect() bool
}
