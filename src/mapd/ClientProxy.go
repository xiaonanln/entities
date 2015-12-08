package mapd

import (
	"fmt"
	"net"
)

type ClientProxy struct {
	MapdConnection
	Pid Pid
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	return &ClientProxy{MapdConnection: NewMapdConnection(conn)}
}

func (self *ClientProxy) SetPid(pid Pid) {
	if self.Pid != 0 {
		panic(fmt.Errorf("SetPid is called twice"))
	}
	self.Pid = pid
}
