package rpc

import (
	"encoding/gob"
	"io"
)

type gobCall struct {
	Eid    string
	Method string
	Args   []interface{}
}

type GobRPCEncoder struct {
	gobEncoder *gob.Encoder
}

func NewGobRPCEncoder(w io.Writer) *GobRPCEncoder {
	return &GobRPCEncoder{gobEncoder: gob.NewEncoder(w)}
}

func (self *GobRPCEncoder) Encode(eid string, method string, arguments []interface{}) error {
	return self.gobEncoder.Encode(gobCall{Eid: eid, Method: method, Args: arguments})
}
