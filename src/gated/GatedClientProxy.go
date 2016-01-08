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

func (self *GatedClientProxy) NewEntity(eid Eid, entityType string) {
}
