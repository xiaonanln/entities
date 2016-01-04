package gated

import (
	"common"
	"net"
)

type GatedClientProxy struct {
	GatedConnection
	ClientId common.ClientId
	Pid      int
}

func NewGatedClientProxy(conn net.Conn) *GatedClientProxy {
	cid := common.NewClientId()
	return &GatedClientProxy{
		GatedConnection: NewGatedConnection(conn),
		ClientId:        cid,
		Pid:             0, // initially no pid
	}
}

func (self *GatedClientProxy) SetPid(pid int) {
	self.Pid = pid
}
