package entitiesd

import (
	. "common"
	"net"
)

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

func (self *EntitiesdClientProxy) NewEntity(clientid ClientId, eid Eid, entityType string) error {
	self.SendCmd(CMD_NEW_ENTITY)
	self.SendCid(clientid)
	self.SendEid(eid)
	return self.SendString(entityType)
}

func (self *EntitiesdClientProxy) DelEntity(clientid ClientId, eid Eid) error {
	self.SendCmd(CMD_DEL_ENTITY)
	self.SendCid(clientid)
	return self.SendEid(eid)
}
