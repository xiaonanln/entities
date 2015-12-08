package main

import (
	"conf"
	"entities"
	"fmt"
	"log"
	. "mapd"
	"net"
	"time"
)

var (
	mapdConfig *conf.MapdConfig
	mapping    map[entities.Eid]Pid
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
		clientOnError(client, err.(error))
	}
	log.Println("Connection closed:", client)
}

func processNextCommand(client *ClientProxy) {
	cmd, err := client.RecvCmd()
	if err != nil {
		panic(err)
	}
	switch cmd {
	case CMD_QUERY:
		processQuery(client)
	case CMD_SET:
		processSet(client)
	case CMD_PID:
		processPid(client)
	case CMD_SYNC_TIME:
		processSyncTime(client)

	}
}

func clientOnError(client *ClientProxy, err error) {
	log.Printf("Error when serving client %s: %s", client, err)
	client.Close()
}

func processPid(client *ClientProxy) {
	pid, err := client.RecvPid()
	if err != nil {
		panic(err)
	}
	client.SetPid(pid)
}

func processQuery(client *ClientProxy) {
	var eid entities.Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}

	pid, ok := mapping[eid]
	if !ok {
		pid = client.Pid
		mapping[eid] = pid
	}
	client.SendPid(pid)
}

func processSet(client *ClientProxy) {
	var eid entities.Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}
	mapping[eid] = client.Pid
	client.SendReplyOk()
}

func processSyncTime(client *ClientProxy) {
	var nano int64 = time.Now().UnixNano()
	client.SendInt64(nano)
}
