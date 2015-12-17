package main

import (
	. "common"
	"conf"
	"fmt"
	"log"
	. "mapd"
	"net"
	"strings"
)

var (
	mapdConfig *conf.MapdConfig
	mapping    map[Eid]Pid = make(map[Eid]Pid)
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
		log.Println("New connection:", conn.RemoteAddr())
		go serveConnection(conn)
	}
}

func serveConnection(conn net.Conn) {
	client := NewClientProxy(conn)
	serveMapdClient(client)
}

func serveMapdClient(client *ClientProxy) {
	defer processClientError(client)
	for {
		processNextCommand(client)
	}
}

func processClientError(client *ClientProxy) {
	if err := recover(); err != nil {
		normalClose := false
		if _err := err.(error); _err != nil {
			errorStr := _err.Error()
			if strings.Contains(errorStr, "EOF") || strings.Contains(errorStr, "closed") {
				// just normal close
				normalClose = true
			}
		}
		if !normalClose {
			log.Printf("ERROR: %s: %T %s", client, err, err)
		}
		client.Close()
	}
	log.Println("Connection closed:", client)
}
