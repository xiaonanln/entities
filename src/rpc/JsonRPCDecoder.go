package rpc

import (
	"encoding/json"
	"io"
)

type jsonCall struct {
	eid    string
	method string
	args   []interface{}
}

type JsonRPCDecoder struct {
	jsonDecoder *json.Decoder
}

func NewJsonRPCDecoder(r io.Reader) *JsonRPCDecoder {
	return &JsonRPCDecoder{jsonDecoder: json.NewDecoder(r)}
}

func (self *JsonRPCDecoder) Decode(eid *string, method *string, args *[]interface{}) error {
	var call jsonCall
	err := self.jsonDecoder.Decode(&call)
	if err != nil {
		return err
	}

	*eid, *method, *args = call.eid, call.method, call.args
	return nil
}
