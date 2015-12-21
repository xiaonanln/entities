package gated

import "net"

type GatedClientProxy struct {
	GatedConnection
}

func NewGatedClientProxy(conn net.Conn) *GatedClientProxy {
	return &GatedClientProxy{GatedConnection: NewGatedConnection(conn)}
}
