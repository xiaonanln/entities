package rpc

import (
	"encoding/gob"
	"io"
)

type GobRPCDecoder struct {
	gobDecoder *gob.Decoder
}

func NewGobRPCDecoder(r io.Reader) *GobRPCDecoder {
	return &GobRPCDecoder{gobDecoder: gob.NewDecoder(r)}
}

func (self *GobRPCDecoder) Decode(eid *string, method *string, args *[]interface{}) error {
	var val gobCall
	err := self.gobDecoder.Decode(&val)
	if err != nil {
		return err
	}
	*eid, *method, *args = val.Eid, val.Method, val.Args
	return nil
}
