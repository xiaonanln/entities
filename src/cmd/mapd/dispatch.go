package main

import (
	. "common"
	"log"
	"mapd"
	"sync"
)

var (
	mapping    = make(map[Eid]int)
	lock       sync.RWMutex                          // only do locking in Caped functions
	rpcClients = make(map[int]*mapd.MapdClientProxy) // TODO: use array instead
)

func AddRPCClient(client *mapd.MapdClientProxy, pid int) {
	lock.Lock()
	rpcClients[pid] = client
	lock.Unlock()

	NotifyAllRegisteredGlobalEntities(client)
}

func DispatchRPC(eid Eid, method string, args []interface{}, fromPid int) {
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

	err := client.RPC(eid, method, args)
	if err != nil {
		handleClientError(client, err) // TODO: use HandleConnectionError instead
	}
}

func DispatchGlobalEntityRegister(eid Eid, entityType string) {
	log.Printf("DISPATCH >>> global entity %s ==> %s", entityType, eid)
	lock.RLock()
	defer lock.RUnlock()

	for _, client := range rpcClients {
		err := client.NotifyRegisterGlobal(eid, entityType)
		if err != nil {
			handleClientError(client, err)
		}
	}
}

func GetMapping(eid Eid, fromPid int) int {
	lock.RLock()
	defer lock.RUnlock()
	return getMapping(eid, fromPid)
}

func getMapping(eid Eid, fromPid int) int {
	pid, ok := mapping[eid]
	if ok {
		return pid
	} else {
		mapping[eid] = fromPid
		return fromPid
	}
}

func SetMapping(eid Eid, pid int) {
	lock.Lock()
	defer lock.Unlock()

	mapping[eid] = pid
	return
}
