package entitiesd

import "net"

type EntitiesdClientProxy struct {
	EntitiesdConnection
	Gid int
}

func NewEntitiesdClientProxy(conn net.Conn) *EntitiesdClientProxy {
	return &EntitiesdClientProxy{
		EntitiesdConnection: NewEntitiesdConnection(conn),
	}
}

func (self *EntitiesdClientProxy) SetGid(gid int) {
	self.Gid = gid
}
