package main

import (
	. "common"
	. "gated"
	"log"
	"sync"
)

type clientSendOp struct {
	opfunc   func(client *GatedClientProxy) error
	clientid ClientId
}

var (
	clientsAccessLock sync.RWMutex
	clients           map[ClientId]*GatedClientProxy
	clientSendOps     = make(chan clientSendOp)
)

func init() {
	go dispatcher()
}

func getClient(clientid ClientId) *GatedClientProxy {
	clientsAccessLock.RLock()
	defer clientsAccessLock.RUnlock()
	return clients[clientid]
}

func setClient(client *GatedClientProxy) {
	clientsAccessLock.Lock()
	defer clientsAccessLock.Unlock()
	clients[client.ClientId] = client
}

func dispatchAddNewClient(client *GatedClientProxy) {
	setClient(client)
}

func dispatchOnClientClose(client *GatedClientProxy) {
	clientsAccessLock.Lock()
	defer clientsAccessLock.Unlock()
	delete(clients, client.ClientId)
}

func dispatchNewEntityToClient(clientid ClientId, eid Eid, entityType string) {
	clientSendOps <- clientSendOp{
		opfunc: func(client *GatedClientProxy) error {
			return nil
		},
		clientid: clientid,
	}

}

func dispatcher() {
	for {
		op := <-clientSendOps
		client := getClient(op.clientid)
		if client == nil {
			// client closed, ignore
			log.Printf("Client %s already disconnected, ignoring op %v", op.clientid, op.opfunc)
			continue
		}

		err := op.opfunc(client)
		if err != nil {
			log.Printf("Client %s dispatch error: %s", client, err)
			continue
		}
	}
}
