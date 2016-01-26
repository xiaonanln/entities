package mapd

import (
	. "common"
	"fmt"
	"log"
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

func (self *MapdClientProxy) NotifyRegisterGlobal(eid Eid, entityType string) error {
	log.Printf(">>> %s: global %s registered to be %s", self, entityType, eid)
	self.SendCmd(CMD_REGISTER_GLOBAL)
	self.SendEid(eid)
	return self.SendString(entityType)
}

func (self *MapdClientProxy) RPC(eid Eid, method string, args []interface{}) error {
	self.SendCmd(CMD_RPC)
	return self.SendRPC(eid, method, args)
}
