package main

import (
	. "common"
	"conf"
	"fmt"
	"log"
	. "mapd"
	"net"
	"strings"
	"sync"
)

var (
	mapdConfig   *conf.MapdConfig
	mapping      map[Eid]Pid = make(map[Eid]Pid)
	waitServices             = sync.WaitGroup{}
)

func main() {
	readConfig()
	startMapdServices()
	waitServices.Wait()
}

func readConfig() {
	mapdConfig = conf.ReadMapdConfig()
	log.Printf("Using mapd config: %v\n", *mapdConfig)
}

func startMapdServices() {
	waitServices.Add(2) // number of services
	go serveCmdService()
	go serveRPCService()
}

func serveCmdService() {
	addr := fmt.Sprintf("%s:%d", mapdConfig.Host, mapdConfig.CmdPort)
	servsock, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("Serving mapd cmd service on", addr, "...")
	for {
		conn, err := servsock.Accept()
		if err != nil {
			panic(err)
		}
		log.Println("New connection:", conn.RemoteAddr())
		go serveCmdConnection(conn)
	}
	waitServices.Done()
}

func serveRPCService() {
	addr := fmt.Sprintf("%s:%d", mapdConfig.Host, mapdConfig.RPCPort)
	servsock, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("Serving mapd RPC service on", addr, "...")
	for {
		conn, err := servsock.Accept()
		if err != nil {
			panic(err)
		}
		log.Println("New connection:", conn.RemoteAddr())
		go serveCmdConnection(conn)
	}
	waitServices.Done()
}

func serveCmdConnection(conn net.Conn) {
	client := NewClientProxy(conn)
	serveCmdMapdClient(client)
}

func serveCmdMapdClient(client *ClientProxy) {
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
