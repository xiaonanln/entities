package entitiesd

import (
	"common"
	"net"
)

// EntitiesdClient for client-server communication
type EntitiesdClient struct {
	EntitiesdConnection
	Pid int
}

func NewEntitiesdClient(conn net.Conn, pid int) *EntitiesdClient {
	return &EntitiesdClient{
		EntitiesdConnection: NewEntitiesdConnection(conn),
		Pid:                 pid,
	}
}

func (self *EntitiesdClient) NewClient(cid common.ClientId) {
	self.SendCmd(CMD_NEW_CLIENT)
	return self.SendCid(cid)
}
