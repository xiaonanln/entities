package rpc

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

type CustomRPCEncoder struct {
	writer io.Writer
}

func NewCustomRPCEncoder(w io.Writer) *CustomRPCEncoder {
	return &CustomRPCEncoder{writer: w}
}

func (self *CustomRPCEncoder) writeString(s string) error {
	lenByte := byte(len(s))
	err := self.writeByte(lenByte)
	if err != nil {
		return err
	}
	return self.writeAll([]byte(s))
}

func (self *CustomRPCEncoder) writeAll(buf []byte) error {
	for len(buf) > 0 {
		n, err := self.writer.Write(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func (self *CustomRPCEncoder) writeByte(b byte) error {
	return self.writeAll([]byte{b})
}

func (self *CustomRPCEncoder) Encode(eid string, method string, arguments []interface{}) error {
	self.writeString(eid)
	self.writeString(method)
	bytes, err := json.Marshal(arguments)
	if err != nil {
		return err
	}

	lengthBytes := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(bytes)))
	self.writeAll(lengthBytes)
	return self.writeAll(bytes)
}
