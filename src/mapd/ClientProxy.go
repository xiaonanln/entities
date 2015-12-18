package mapd

import (
	"common"
	"fmt"
	"net"
	"rpc"
)

type ClientProxy struct {
	MapdConnection
	Pid        Pid
	rpcEncoder rpc.RPCEncoder
	rpcDecoder rpc.RPCDecoder
}

func NewClientProxy(conn net.Conn) *ClientProxy {
	cp := &ClientProxy{MapdConnection: NewMapdConnection(conn)}
	cp.rpcEncoder = rpc.NewCustomRPCEncoder(cp.MapdConnection)
	cp.rpcDecoder = rpc.NewCustomRPCDecoder(cp.MapdConnection)
	return cp
}

func (self *ClientProxy) SetPid(pid Pid) {
	if self.Pid != 0 {
		panic(fmt.Errorf("SetPid is called twice"))
	}
	self.Pid = pid
}

func (self *ClientProxy) OnRPC(eid common.Eid, method string, args []interface{}) error {
	return self.SendRPC(eid, method, args)
}
