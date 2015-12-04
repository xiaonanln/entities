package rpc

type RPCEncoder interface {
	Encode(eid string, method string, args []interface{}) error
}
type RPCDecoder interface {
	Decode(eid *string, method *string, args *[]interface{}) error
}
