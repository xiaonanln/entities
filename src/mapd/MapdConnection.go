package mapd

import (
	"common"
	"net"
)

const (
	MAPD_OP_QUERY  = 'R'
	MAPD_OP_CREATE = 'C'
)

type MapdConnection struct {
	common.BinaryConnection
}

func NewMapdConnection(conn net.Conn) MapdConnection {
	binaryConn := common.NewBinaryConnection(conn)
	return MapdConnection{BinaryConnection: binaryConn}
}
