package rpc

import (
	"encoding/json"
	"io"
)

type JsonRPCEncoder struct {
	jsonEncoder *json.Encoder
}

func NewJsonRPCEncoder(w io.Writer) *JsonRPCEncoder {
	return &JsonRPCEncoder{jsonEncoder: json.NewEncoder(w)}
}

func (self *JsonRPCEncoder) Encode(eid string, method string, arguments []interface{}) error {
	return self.jsonEncoder.Encode(jsonCall{eid: eid, method: method, args: arguments})
}
