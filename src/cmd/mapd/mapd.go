package main

import (
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
	waitServices = sync.WaitGroup{}
)

func main() {
	log.SetPrefix("mapd ")

	readConfig()
	startMapdServices()
	waitServices.Wait()
}

func readConfig() {
	mapdConfig = conf.GetMapdConfig()
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
		go serveRPCConnection(conn)
	}
	waitServices.Done()
}

func serveRPCConnection(conn net.Conn) {
	client := NewClientProxy(conn)
	pid, err := client.RecvPid()
	if err != nil {
		handleClientError(client, err)
		return
	}
	log.Printf("Received pid of %s: %v", client, pid)
	client.SetPid(pid)
	AddRPCClient(client, pid)
}

func serveCmdConnection(conn net.Conn) {
	client := NewClientProxy(conn)
	serveCmdMapdClient(client)
}

func serveCmdMapdClient(client *MapdClientProxy) {
	defer serveCmdMapdClientDone(client)
	for {
		processNextCommand(client)
	}
}

func serveCmdMapdClientDone(client *MapdClientProxy) {
	if err := recover(); err != nil {
		handleClientError(client, err)
	}
	log.Println("Connection closed:", client)
}

func handleClientError(client *MapdClientProxy, err interface{}) {
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
