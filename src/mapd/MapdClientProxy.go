package mapd

import (
	"common"
	"fmt"
	"net"
	"rpc"
)

type MapdClientProxy struct {
	MapdConnection
	Pid        int
	rpcEncoder rpc.RPCEncoder
	rpcDecoder rpc.RPCDecoder
}

func NewClientProxy(conn net.Conn) *MapdClientProxy {
	cp := &MapdClientProxy{MapdConnection: NewMapdConnection(conn)}
	cp.rpcEncoder = rpc.NewCustomRPCEncoder(cp.MapdConnection)
	cp.rpcDecoder = rpc.NewCustomRPCDecoder(cp.MapdConnection)
	return cp
}

func (self *MapdClientProxy) SetPid(pid int) {
	if self.Pid != 0 {
		panic(fmt.Errorf("SetPid is called twice"))
	}
	self.Pid = pid
}

func (self *MapdClientProxy) OnRPC(eid common.Eid, method string, args []interface{}) error {
	return self.SendRPC(eid, method, args)
}
