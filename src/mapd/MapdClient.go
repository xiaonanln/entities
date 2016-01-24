package mapd

import (
	. "common"
	"net"
)

type MapdClient struct {
	MapdConnection
	pid int
}

func NewMapdClient(conn net.Conn) *MapdClient {
	client := &MapdClient{MapdConnection: NewMapdConnection(conn)}
	return client
}

func (self *MapdClient) SetPid(pid int) error {
	self.pid = pid

	self.SendCmd(CMD_PID)
	err := self.SendPid(pid)
	if err != nil {
		return err
	}
	return self.RecvReplyOk()
}

func (self *MapdClient) SetMapping(eid Eid) error {
	var err error
	self.SendCmd(CMD_SET)
	err = self.SendEid(eid)
	if err != nil {
		return err
	}

	return self.RecvReplyOk()
}

func (self *MapdClient) RegisterGlobal(eid Eid, entityType string) (bool, error) {
	var err error
	self.SendCmd(CMD_REGISTER_GLOBAL)
	self.SendEid(eid)
	err = self.SendString(entityType)
	if err != nil {
		return false, err
	}
	ret, err := self.RecvByte()
	if err != nil {
		return false, err
	}
	return ret == REPLY_OK, nil
}

func (self *MapdClient) RPC(eid Eid, method string, args []interface{}) error {
	self.SendCmd(CMD_RPC)
	err := self.SendRPC(eid, method, args)
	if err != nil {
		return err
	}
	return self.RecvReplyOk()
}
