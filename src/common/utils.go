package common

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func ConnectTCP(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	return conn, err
}

func HandleConnectionError(conn interface{}, err interface{}) {
	normalClose := false
	if _err := err.(error); _err != nil {
		errorStr := _err.Error()
		if strings.Contains(errorStr, "EOF") || strings.Contains(errorStr, "closed") {
			// just normal close
			normalClose = true
		}
	}
	if !normalClose {
		log.Printf("ERROR: %s: %T %s", conn, err, err)
	} else {
		log.Println("Connection closed:", conn)
	}
}
