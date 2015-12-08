package common

import "net"

type BinaryConnection struct {
	Connection
}

func NewBinaryConnection(conn net.Conn) BinaryConnection {
	return BinaryConnection{Connection{conn}}
}
