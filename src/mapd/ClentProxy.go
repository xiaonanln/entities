package mapd

import (
	"common"
	"net"
)

type ClientProxy struct {
	common.BinaryConnection
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{common.NewBinaryConnection(conn)}
}
