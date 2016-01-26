package mapd

import (
	. "common"
	"fmt"
	"net"
	"rpc"
)

type MapdConnection struct {
	BinaryConnection
	rpcEncoder rpc.RPCEncoder
	rpcDecoder rpc.RPCDecoder
}

func NewMapdConnection(conn net.Conn) MapdConnection {
	binaryConn := NewBinaryConnection(conn)
	mapdCon := MapdConnection{BinaryConnection: binaryConn}
	mapdCon.rpcEncoder = rpc.NewCustomRPCEncoder(mapdCon)
	mapdCon.rpcDecoder = rpc.NewCustomRPCDecoder(mapdCon)
	return mapdCon
}

func (self *MapdConnection) RecvCmd() (byte, error) {
	b, err := self.RecvByte()
	return b, err
}

func (self *MapdConnection) SendCmd(cmd byte) error {
	return self.SendByte(cmd)
}

func (self *MapdConnection) SendReplyOk() error {
	return self.SendByte(REPLY_OK)
}

func (self *MapdConnection) RecvReplyOk() error {
	b, err := self.RecvByte()
	if err != nil {
		return err
	}
	if b != REPLY_OK {
		return fmt.Errorf("expect REPLY_OK but received %v", b)
	}
	return nil
}

func (self *MapdConnection) SendRPC(eid Eid, method string, args []interface{}) error {
	return self.rpcEncoder.Encode(string(eid), method, args)
}

func (self *MapdConnection) RecvRPC(eid *Eid, method *string, args *[]interface{}) error {
	return self.rpcDecoder.Decode((*string)(eid), method, args)
}

func (self *MapdConnection) SendPid(pid int) error {
	return self.SendUint16(uint16(pid))
}

func (self *MapdConnection) RecvPid() (int, error) {
	v, err := self.RecvUint16()
	return int(v), err
}
