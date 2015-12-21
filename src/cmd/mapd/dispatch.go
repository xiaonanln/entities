package main

import (
	"common"
	"log"
	"mapd"
	"sync"
)

var (
	mapping    = make(map[common.Eid]int)
	lock       sync.RWMutex // only do locking in Caped functions
	rpcClients = make(map[int]*mapd.MapdClientProxy)
)

func AddRPCClient(client *mapd.MapdClientProxy, pid int) {
	lock.Lock()
	defer lock.Unlock()

	rpcClients[pid] = client
}

func DispatchRPC(eid common.Eid, method string, args []interface{}, fromPid int) {
	log.Printf("DISPATCH >>> %s.%s%v", eid, method, args)
	lock.RLock()
	defer lock.RUnlock()

	pid := getMapping(eid, fromPid)
	// send to pid now
	client, ok := rpcClients[pid]
	if !ok {
		// client not found, which should not happen
		log.Println("RPC client for pid", pid, "is not found, rpc failed!")
		return
	}

	err := client.OnRPC(eid, method, args)
	if err != nil {
		handleClientError(client, err)
	}
}

func GetMapping(eid common.Eid, fromPid int) int {
	lock.RLock()
	defer lock.RUnlock()
	return getMapping(eid, fromPid)
}

func getMapping(eid common.Eid, fromPid int) int {
	pid, ok := mapping[eid]
	if ok {
		return pid
	} else {
		mapping[eid] = fromPid
		return fromPid
	}
}

func SetMapping(eid common.Eid, pid int) {
	lock.Lock()
	defer lock.Unlock()

	mapping[eid] = pid
	return
}
