package rpc

import (
	"encoding/json"
	"io"
)

type CustomRPCDecoder struct {
	reader      io.Reader
	jsonDecoder *json.Decoder
}

func NewCustomRPCDecoder(r io.Reader) *CustomRPCDecoder {
	return &CustomRPCDecoder{reader: r, jsonDecoder: json.NewDecoder(r)}
}

func (self *CustomRPCDecoder) Decode(eid *string, method *string, args *[]interface{}) error {
	self.readString(eid)
	self.readString(method)
	return self.jsonDecoder.Decode(args)
}

func (self *CustomRPCDecoder) readString(s *string) error {
	lenByte, err := self.readByte()
	if err != nil {
		return err
	}

	strLen := int(lenByte)
	strbuf := make([]byte, strLen, strLen)
	err = self.readAll(strbuf)
	if err != nil {
		return err
	}
	*s = string(strbuf)
	return nil
}

func (self *CustomRPCDecoder) readByte() (byte, error) {
	buf := []byte{0}
	err := self.readAll(buf)
	return buf[0], err
}

func (self *CustomRPCDecoder) readAll(p []byte) error {
	for len(p) > 0 {
		n, err := self.reader.Read(p)
		if err != nil {
			return err
		}
		p = p[n:]
	}
	return nil
}
