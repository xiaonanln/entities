package entitiesd

import (
	. "common"
	"net"
	"rpc"
)

type EntitiesdConnection struct {
	BinaryConnection
	rpcEncoder rpc.RPCEncoder
	rpcDecoder rpc.RPCDecoder
}

func NewEntitiesdConnection(conn net.Conn) EntitiesdConnection {
	binaryConn := NewBinaryConnection(conn)
	gatedCon := EntitiesdConnection{BinaryConnection: binaryConn}
	gatedCon.rpcEncoder = rpc.NewCustomRPCEncoder(gatedCon)
	gatedCon.rpcDecoder = rpc.NewCustomRPCDecoder(gatedCon)
	return gatedCon
}

func (self EntitiesdConnection) RecvCmd() (byte, error) {
	b, err := self.RecvByte()
	return b, err
}

func (self EntitiesdConnection) SendCmd(cmd byte) error {
	return self.SendByte(cmd)
}

func (self EntitiesdConnection) SendRPC(eid Eid, method string, args []interface{}) error {
	return self.rpcEncoder.Encode(string(eid), method, args)
}

func (self EntitiesdConnection) RecvRPC(eid *Eid, method *string, args *[]interface{}) error {
	return self.rpcDecoder.Decode((*string)(eid), method, args)
}

func (self EntitiesdConnection) SendGid(pid int) error {
	return self.SendUint16(uint16(pid))
}

func (self EntitiesdConnection) RecvGid() (int, error) {
	v, err := self.RecvUint16()
	return int(v), err
}

func (self EntitiesdConnection) SendCid(cid ClientId) error {
	return self.SendFixedLengthString(string(cid))
}

func (self EntitiesdConnection) RecvCid(cid *ClientId) error {
	return self.RecvFixedLengthString(EID_LENGTH, (*string)(cid))
}

func (self EntitiesdConnection) RPC(eid Eid, method string, args []interface{}) error {
	self.SendCmd(CMD_RPC)
	return self.SendRPC(eid, method, args)
}
