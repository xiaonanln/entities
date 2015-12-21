package gated

import (
	. "common"
	"net"
	"rpc"
)

type GatedConnection struct {
	BinaryConnection
	rpcEncoder rpc.RPCEncoder
	rpcDecoder rpc.RPCDecoder
}

func NewGatedConnection(conn net.Conn) GatedConnection {
	binaryConn := NewBinaryConnection(conn)
	gatedCon := GatedConnection{BinaryConnection: binaryConn}
	gatedCon.rpcEncoder = rpc.NewCustomRPCEncoder(gatedCon)
	gatedCon.rpcDecoder = rpc.NewCustomRPCDecoder(gatedCon)
	return gatedCon
}

func (self GatedConnection) RecvCmd() (byte, error) {
	b, err := self.RecvByte()
	return b, err
}

func (self GatedConnection) SendCmd(cmd byte) error {
	return self.SendByte(cmd)
}

func (self GatedConnection) SendRPC(eid Eid, method string, args []interface{}) error {
	return self.rpcEncoder.Encode(string(eid), method, args)
}

func (self GatedConnection) RecvRPC(eid *Eid, method *string, args *[]interface{}) error {
	return self.rpcDecoder.Decode((*string)(eid), method, args)
}

func (self GatedConnection) SendGid(pid int) error {
	return self.SendUint16(uint16(pid))
}

func (self GatedConnection) RecvGid() (int, error) {
	v, err := self.RecvUint16()
	return int(v), err
}

func (self GatedConnection) RPC(eid Eid, method string, args []interface{}) error {
	self.SendCmd(CMD_RPC)
	return self.SendRPC(eid, method, args)
}
