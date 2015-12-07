package main

import (
	"conf"
	"fmt"
	"log"
	"net"
)

var (
	mapdConfig *conf.MapdConfig
)

func main() {
	readConfig()
	runMapd()
}

func readConfig() {
	mapdConfig = conf.ReadMapdConfig()
	log.Printf("Using mapd config: %v\n", *mapdConfig)
}

func runMapd() {
	addr := fmt.Sprintf("%s:%d", mapdConfig.Host, mapdConfig.Port)
	servsock, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("Serving mapd on", addr, "...")
	for {
		conn, err := servsock.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

}
