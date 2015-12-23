package entitiesd

import "net"

type EntitiesdClientProxy struct {
	EntitiesdConnection
}

func NewEntitiesdClientProxy(conn net.Conn) *EntitiesdClientProxy {
	return &EntitiesdClientProxy{
		EntitiesdConnection: NewEntitiesdConnection(conn),
	}
}
