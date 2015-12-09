package mapd

import (
	"entities"
	"net"
)

type MapdClient struct {
	MapdConnection
	pid Pid
}

func NewMapdClient(conn net.Conn) *MapdClient {
	return &MapdClient{MapdConnection: NewMapdConnection(conn)}
}

func (self *MapdClient) SetPid(pid Pid) error {
	self.pid = pid

	self.SendCmd(CMD_PID)
	err := self.SendPid(pid)
	if err != nil {
		return err
	}
	return self.RecvReplyOk()
}

func (self *MapdClient) SetMapping(eid entities.Eid) error {
	var err error
	self.SendCmd(CMD_SET)
	err = self.SendEid(eid)
	if err != nil {
		return err
	}

	return self.RecvReplyOk()
}
