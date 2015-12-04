package rpc

import (
	"encoding/json"
	"io"
)

type CustomRPCEncoder struct {
	writer      io.Writer
	jsonEncoder *json.Encoder
}

func NewCustomRPCEncoder(w io.Writer) *CustomRPCEncoder {
	return &CustomRPCEncoder{writer: w, jsonEncoder: json.NewEncoder(w)}
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
	buf := []byte{b}
	return self.writeAll(buf)
}

func (self *CustomRPCEncoder) Encode(eid string, method string, arguments []interface{}) error {
	self.writeString(eid)
	self.writeString(method)
	return self.jsonEncoder.Encode(arguments)
}
