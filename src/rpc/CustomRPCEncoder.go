package rpc

import (
	"encoding/json"
	"io"
)

const (
	MAX_EID_LENGTH    = 16
	MAX_METHOD_LENGTH = 32
)

type CustomRPCEncoder struct {
	writer      io.Writer
	jsonEncoder *json.Encoder
}

func NewCustomRPCEncoder(w io.Writer) *CustomRPCEncoder {
	return &CustomRPCEncoder{writer: w, jsonEncoder: json.NewEncoder(w)}
}

func (self *CustomRPCEncoder) writeFixedLenString(s string, fixedLen int) error {
	var err error
	bytes := []byte(s)
	if len(bytes) >= fixedLen {
		return self.writeAll(bytes[:fixedLen])
	} else {
		err = self.writeAll(bytes)
		if err != nil {
			return nil
		}
		return self.writeByte(' ', fixedLen-len(bytes))
	}
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

func (self *CustomRPCEncoder) writeByte(b byte, count int) error {
	bytes := make([]byte, count, count)
	for i := 0; i < count; i++ {
		bytes[i] = b
	}
	return self.writeAll(bytes)
}

func (self *CustomRPCEncoder) Encode(eid string, method string, arguments []interface{}) error {
	err := self.writeFixedLenString(eid, MAX_EID_LENGTH)
	if err != nil {
		return err
	}

	err = self.writeFixedLenString(method, MAX_METHOD_LENGTH)
	if err != nil {
		return err
	}
	return self.jsonEncoder.Encode(arguments)
}
