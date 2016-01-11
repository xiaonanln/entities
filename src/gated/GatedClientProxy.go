package gated

import (
	. "common"
	"net"
)

type GatedClientProxy struct {
	GatedConnection
	ClientId ClientId
	Pid      int
}

func NewGatedClientProxy(conn net.Conn) *GatedClientProxy {
	cid := NewClientId()
	return &GatedClientProxy{
		GatedConnection: NewGatedConnection(conn),
		ClientId:        cid,
		Pid:             0, // initially no pid
	}
}

func (self *GatedClientProxy) SetPid(pid int) {
	self.Pid = pid
}

func (self *GatedClientProxy) NewEntity(eid Eid, entityType string) error {
	self.SendCmd(CMD_NEW_ENTITY)
	self.SendEid(eid)
	return self.SendString(entityType)
}

func (self *GatedClientProxy) DelEntity(eid Eid) error {
	self.SendCmd(CMD_DEL_ENTITY)
	return self.SendEid(eid)
}
