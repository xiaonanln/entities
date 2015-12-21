package gated

import "net"

type GatedClient struct {
	GatedConnection
}

func NewGatedClient(conn net.Conn) *GatedClient {
	return &GatedClient{GatedConnection: NewGatedConnection(conn)}
}
