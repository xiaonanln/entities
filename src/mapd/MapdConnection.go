package mapd

import (
	"common"
	"entities"
	"fmt"
	"net"
)

const (
	MAPD_OP_QUERY  = 'R'
	MAPD_OP_CREATE = 'C'
)

type Pid uint16

type MapdConnection struct {
	common.BinaryConnection
}

func NewMapdConnection(conn net.Conn) MapdConnection {
	binaryConn := common.NewBinaryConnection(conn)
	return MapdConnection{BinaryConnection: binaryConn}
}

func (self MapdConnection) RecvCmd() (byte, error) {
	b, err := self.RecvByte()
	return b, err
}

func (self MapdConnection) SendCmd(cmd byte) error {
	return self.SendByte(cmd)
}

func (self MapdConnection) RecvEid(eid *entities.Eid) error {
	err := self.RecvFixedLengthString(entities.EID_LENGTH, (*string)(eid))
	return err
}

func (self MapdConnection) SendEid(eid entities.Eid) error {
	return self.SendFixedLengthString(string(eid))
}

func (self MapdConnection) SendPid(pid Pid) error {
	return self.SendUint16(uint16(pid))
}

func (self MapdConnection) RecvPid() (Pid, error) {
	v, err := self.RecvUint16()
	return Pid(v), err
}

func (self MapdConnection) SendReplyOk() error {
	return self.SendByte(REPLY_OK)
}

func (self MapdConnection) RecvReplyOk() error {
	b, err := self.RecvByte()
	if err != nil {
		return err
	}
	if b != REPLY_OK {
		return fmt.Errorf("expect REPLY_OK but received %v", b)
	}
	return nil
}
