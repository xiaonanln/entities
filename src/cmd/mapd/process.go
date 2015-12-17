package main

import (
	. "common"
	"log"
	. "mapd"
	"time"
)

func processNextCommand(client *ClientProxy) {
	cmd, err := client.RecvCmd()
	if err != nil {
		panic(err)
	}
	log.Println("CMD:", cmd)
	switch cmd {
	case CMD_RPC:
		processRPC(client)
	case CMD_QUERY:
		processQuery(client)
	case CMD_SET:
		processSet(client)
	case CMD_PID:
		processPid(client)
	// case CMD_SYNC_TIME:
	// 	processSyncTime(client)
	case CMD_LOCK_EID:
		processLockEid(client)
	}
}

func processPid(client *ClientProxy) {
	pid, err := client.RecvPid()
	if err != nil {
		panic(err)
	}
	client.SetPid(pid)
	client.SendReplyOk()
}

func processQuery(client *ClientProxy) {
	var eid Eid
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
	var eid Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}
	mapping[eid] = client.Pid
	client.SendReplyOk()
	log.Printf("SET %s => %d", eid, client.Pid)
}

func processSyncTime(client *ClientProxy) {
	var nano int64 = time.Now().UnixNano()
	client.SendInt64(nano)
}

func processLockEid(client *ClientProxy) {
	var eid Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}
	// TODO: lock Eid, cache post-coming calls for a period of time
	// send lock ok when ready
}

func processRPC(client *ClientProxy) {
	var eid Eid
	var method string
	var args []interface{}
	err := client.RecvRPC(&eid, &method, &args)
	if err != nil {
		panic(err)
	}
	log.Printf("RPC: %s.%s(%v)", eid, method, args)
	// send back ok
	client.SendReplyOk()
}
