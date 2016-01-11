package gated

import (
	. "common"
	"log"
	"net"
)

type GatedClient struct {
	GatedConnection
}

func NewGatedClient(conn net.Conn) *GatedClient {
	return &GatedClient{GatedConnection: NewGatedConnection(conn)}
}

func (self *GatedClient) Call(eid Eid, method string, args ...interface{}) error {
	log.Printf("RPC >>> %s.%s%v", eid, method, args)
	return self.RPC(eid, method, args)
}
