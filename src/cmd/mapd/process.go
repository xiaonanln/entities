package main

import (
	. "common"
	"log"
	. "mapd"
	"time"
)

func processNextCommand(client *MapdClientProxy) {
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
	case CMD_REGISTER_GLOBAL:
		processRegisterGlobal(client)
	}
}

func processPid(client *MapdClientProxy) {
	pid, err := client.RecvPid()
	if err != nil {
		panic(err)
	}
	client.SetPid(pid)
	client.SendReplyOk()
}

func processQuery(client *MapdClientProxy) {
	var eid Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}

	pid := GetMapping(eid, client.Pid)
	client.SendPid(pid)
}

func processSet(client *MapdClientProxy) {
	var eid Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}
	SetMapping(eid, client.Pid)
	client.SendReplyOk()
	log.Printf("SET %s => %d", eid, client.Pid)
}

func processSyncTime(client *MapdClientProxy) {
	var nano int64 = time.Now().UnixNano()
	client.SendInt64(nano)
}

func processLockEid(client *MapdClientProxy) {
	var eid Eid
	err := client.RecvEid(&eid)
	if err != nil {
		panic(err)
	}
	// TODO: lock Eid, cache post-coming calls for a period of time
	// send lock ok when ready
}

func processRegisterGlobal(client *MapdClientProxy) {
	var eid Eid
	var entityType string
	client.RecvEid(&eid)
	if err := client.RecvString(&entityType); err != nil {
		panic(err)
	}

	if registerGlobalEntity(eid, entityType) {
		client.SendReplyOk()
	} else {
		client.SendByte(REPLY_REGISTER_FAIL)
	}
}

func processRPC(client *MapdClientProxy) {
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

	DispatchRPC(eid, method, args, client.Pid)
}
