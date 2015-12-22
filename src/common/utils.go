package common

import (
	"fmt"
	"net"
)

func ConnectTCP(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	return conn, err
}
