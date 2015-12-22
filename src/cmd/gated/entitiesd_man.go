package main

import (
	"common"
	"conf"
	"entitiesd"
	"gated"
	"log"
	"math/rand"
	"time"
)

var (
	entitiesdClients []*entitiesd.EntitiesdClient
	entitiesdConfigs []conf.EntitiesdConfig
)

func init() {
	go maintainEntitiesdConnections()
}

func maintainEntitiesdConnections() {
	config := conf.GetEntitiesConfig()
	entitiesdConfigs = config.Entitiesd // copy config
	log.Printf("found %d entitiesd", len(entitiesdConfigs))
	entitiesdClients = make([]*entitiesd.EntitiesdClient, len(entitiesdConfigs))

	for {
		for i, _ := range entitiesdConfigs {
			if entitiesdClients[i] == nil {
				connectEntitiesd(i)
			}
		}

		time.Sleep(time.Second)
	}
}

func connectEntitiesd(pid int) {
	config := entitiesdConfigs[pid]
	log.Println("Connecting to entitiesd[%d] @ %s:%d", pid, config.Host, config.Port)

	conn, err := common.ConnectTCP(config.Host, config.Port)
	if err != nil {
		log.Printf("Connect error: %s", err)
		return
	}

	client := entitiesd.NewEntitiesdClient(conn, pid)
	entitiesdClients[pid] = client
	log.Printf("Connected successfully")
}

func chooseRandomEntitiesd() *entitiesd.EntitiesdClient {
	if len(entitiesdClients) == 0 {
		return nil
	}
	return entitiesdClients[rand.Intn(len(entitiesdClients))]
}

func onClientConnect(client *gated.GatedClientProxy) {
	// client connected, choose random entitiesd and tell it
	entitiesd := chooseRandomEntitiesd()
	if entitiesd == nil {
		// connect to entitiesd failed, tell the gate client to shutdown
		log.Printf("Found nil entitiesd, disconnecting client %s", client)
		client.Close()
		return
	}

	client.SetPid(entitiesd.Pid)
	entitiesd.NewClient(client.Cid)
}

func onClientCallRPC(client *gated.GatedClientProxy, eid common.Eid, method string, args []interface{}) {
	entitiesd := entitiesdClients[client.Pid]
	if entitiesd == nil {
		log.Printf("Found nil entitiesd %d, RPC dropped", client.Pid)
		return
	}

	entitiesd.RPC(eid, method, args)
}
