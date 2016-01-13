package main

import (
	. "common"
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

	conn, err := ConnectTCP(config.Host, config.Port)
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
	go serviceEntitiesdClient(client)

	log.Printf("Connected successfully")
}

func chooseRandomEntitiesd() (int, *entitiesd.EntitiesdClient) {
	if len(entitiesdClients) == 0 {
		return -1, nil
	}
	i := rand.Intn(len(entitiesdClients))
	return i + 1, entitiesdClients[i]
}

func letClientConnectEntitiesd(client *gated.GatedClientProxy) error {
	// client connected, choose random entitiesd and tell it
	pid, entitiesd := chooseRandomEntitiesd()
	if entitiesd == nil {
		// connect to entitiesd failed, tell the gate client to shutdown
		log.Printf("Found nil entitiesd, disconnecting client %s", client)
		return errors.New("entitiesd is nil")
	}

	client.SetPid(entitiesd.Pid)
	err := entitiesd.NewClient(client.ClientId)
	if err != nil {
		entitiesdClientOnError(pid, entitiesd, err)
	}
	return err
}

func letClientRPCtoEntitiesd(client *gated.GatedClientProxy, eid Eid, method string, args []interface{}) error {
	log.Printf("RPC from client %s: %s.%s%v", client, eid, method, args)
	entitiesd := entitiesdClients[client.Pid-1]
	if entitiesd == nil {
		log.Printf("Found nil entitiesd %d, RPC dropped", client.Pid)
		return errors.New("entitiesd is nil")
	}

	err := entitiesd.RPC(client.ClientId, eid, method, args)
	if err != nil {
		entitiesdClientOnError(client.Pid, entitiesd, err)
	}
	return err
}

func letClientDisconnectEntitiesd(client *gated.GatedClientProxy) error {
	if client.Pid <= 0 {
		return nil
	}
	entitiesd := entitiesdClients[client.Pid-1]
	if entitiesd == nil {
		return nil
	}
	if entitiesd != nil {
		if err := entitiesd.DelClient(client.ClientId); err != nil {
			entitiesdClientOnError(client.Pid, entitiesd, err)
		}
	}
	return nil
}

func entitiesdClientOnError(pid int, entitiesd *entitiesd.EntitiesdClient, err error) {
	entitiesdClients[pid-1] = nil
}
