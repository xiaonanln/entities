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

func (self BinaryConnection) RecvFixedLengthString(len int, pstr *string) error {
	buf := make([]byte, len)
	err := self.RecvAll(buf)
	if err != nil {
		return err
	}
	*pstr = string(buf)
	return nil
}

func (self BinaryConnection) SendFixedLengthString(s string) error {
	return self.SendAll([]byte(s))
}

func (self BinaryConnection) RecvUint16() (uint16, error) {
	buf := []byte{0, 0}
	err := self.RecvAll(buf)
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) + (uint16(buf[1]) << 8), nil
}

func (self BinaryConnection) SendUint16(val uint16) error {
	buf := []byte{byte(val), byte(val >> 8)}
	return self.SendAll(buf)
}

func (self BinaryConnection) SendInt64(val int64) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(val))
	return self.SendAll(bytes)
}
