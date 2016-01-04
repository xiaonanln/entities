package main

import (
	"common"
	"conf"
	"entitiesd"
	"errors"
	"gated"
	"log"
	"math/rand"
	"setup"
	"time"
)

var (
	entitiesdClients []*entitiesd.EntitiesdClient
)

func init() {
	if setup.IsGated() {
		go maintainEntitiesdConnections()
	}
}

func maintainEntitiesdConnections() {
	log.Printf("found %d entitiesd", conf.GetEntitiesdCount())
	entitiesdClients = make([]*entitiesd.EntitiesdClient, conf.GetEntitiesdCount())

	for {
		for i, client := range entitiesdClients {
			if client == nil {
				connectEntitiesd(i + 1)
			}
		}

		time.Sleep(time.Second)
	}
}

func connectEntitiesd(pid int) {
	config := conf.GetEntitiesdConfig(pid)
	log.Printf("Connecting to entitiesd[%d] @ %s:%d", pid, config.Host, config.Port)

	conn, err := common.ConnectTCP(config.Host, config.Port)
	if err != nil {
		log.Printf("Connect error: %s", err)
		return
	}

	client := entitiesd.NewEntitiesdClient(conn, pid)
	err = client.SendGid(gid)
	if err != nil {
		log.Printf("Send gid error: %s", err)
		return
	}

	entitiesdClients[pid-1] = client
	log.Printf("Connected successfully")
}

func chooseRandomEntitiesd() *entitiesd.EntitiesdClient {
	if len(entitiesdClients) == 0 {
		return nil
	}
	return entitiesdClients[rand.Intn(len(entitiesdClients))]
}

func onClientConnect(client *gated.GatedClientProxy) error {
	// client connected, choose random entitiesd and tell it
	entitiesd := chooseRandomEntitiesd()
	if entitiesd == nil {
		// connect to entitiesd failed, tell the gate client to shutdown
		log.Printf("Found nil entitiesd, disconnecting client %s", client)
		return errors.New("entitiesd is nil")
	}

	client.SetPid(entitiesd.Pid)
	return entitiesd.NewClient(client.ClientId)
}

func onClientCallRPC(client *gated.GatedClientProxy, eid common.Eid, method string, args []interface{}) {
	entitiesd := entitiesdClients[client.Pid]
	if entitiesd == nil {
		log.Printf("Found nil entitiesd %d, RPC dropped", client.Pid)
		return
	}

	entitiesd.RPC(client.ClientId, eid, method, args)
}
