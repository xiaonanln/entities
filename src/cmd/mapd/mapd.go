package main

import (
	"conf"
	"entities"
	"fmt"
	"log"
	. "mapd"
	"net"
)

var (
	mapdConfig *conf.MapdConfig
	mapping    map[entities.Eid]Pid = make(map[entities.Eid]Pid)
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
		log.Printf("ERROR: %s: %s", client, err)
		client.Close()
	}
	log.Println("Connection closed:", client)
}
