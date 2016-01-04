package common

import (
	"encoding/binary"
	"net"
)

type BinaryConnection struct {
	Connection
}

func NewBinaryConnection(conn net.Conn) BinaryConnection {
	return BinaryConnection{Connection{conn}}
}

func (self *BinaryConnection) RecvFixedLengthString(len int, pstr *string) error {
	buf := make([]byte, len)
	err := self.RecvAll(buf)
	if err != nil {
		return err
	}
	*pstr = string(buf)
	return nil
}

func (self *BinaryConnection) SendFixedLengthString(s string) error {
	return self.SendAll([]byte(s))
}

func (self *BinaryConnection) RecvUint16() (uint16, error) {
	buf := []byte{0, 0}
	err := self.RecvAll(buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) + (uint16(buf[1]) << 8), nil
}

func (self *BinaryConnection) SendUint16(val uint16) error {
	buf := []byte{byte(val), byte(val >> 8)}
	return self.SendAll(buf)
}

func (self *BinaryConnection) SendInt64(val int64) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(val))
	return self.SendAll(bytes)
}

func (self *BinaryConnection) RecvEid(eid *Eid) error {
	err := self.RecvFixedLengthString(EID_LENGTH, (*string)(eid))
	return err
}

func (self *BinaryConnection) SendEid(eid Eid) error {
	return self.SendFixedLengthString(string(eid))
}

func (self *BinaryConnection) SendString(s string) error {
	return self.SendByteSlice([]byte(s))
}

func (self *BinaryConnection) RecvString(s *string) error {
	var buf []byte
	err := self.RecvByteSlice(&buf)
	if err != nil {
		return err
	}
	*s = string(buf)
	return nil
}

func (self *BinaryConnection) SendByteSlice(a []byte) error {
	self.SendUint16(uint16(len(a)))
	return self.SendAll(a)
}

func (self *BinaryConnection) RecvByteSlice(a *[]byte) error {
	alen, err := self.RecvUint16()
	if err != nil {
		return err
	}
	buf := make([]byte, alen)
	*a = buf
	return self.RecvAll(buf)
}
