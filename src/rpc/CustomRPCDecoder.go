package rpc

import (
	"encoding/json"
	"io"
	"strings"
)

type CustomRPCDecoder struct {
	reader      io.Reader
	jsonDecoder *json.Decoder
}

func NewCustomRPCDecoder(r io.Reader) *CustomRPCDecoder {
	return &CustomRPCDecoder{reader: r, jsonDecoder: json.NewDecoder(r)}
}

func (self *CustomRPCDecoder) Decode(eid *string, method *string, args *[]interface{}) error {
	err := self.readFixedLenString(eid, MAX_EID_LENGTH)
	if err != nil {
		return nil
	}
	err = self.readFixedLenString(method, MAX_METHOD_LENGTH)
	if err != nil {
		return nil
	}

	return self.jsonDecoder.Decode(args)
}

func (self *CustomRPCDecoder) readFixedLenString(s *string, fixedLen int) error {
	buf := make([]byte, fixedLen, fixedLen)
	err := self.readAll(buf)
	if err != nil {
		return err
	}
	*s = strings.TrimSpace(string(buf))
	return nil
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
