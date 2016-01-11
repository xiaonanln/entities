package main

import (
	. "common"
	"entitiesd"
	"log"
)

func serviceEntitiesdClient(client *entitiesd.EntitiesdClient) {
	for {
		var err error
		cmd, err := client.RecvCmd()
		if err != nil {
			HandleConnectionError(client, err)
			break
		}

		switch cmd {
		case entitiesd.CMD_NEW_ENTITY:
			err = handleNewEntity(client)
		case entitiesd.CMD_DEL_ENTITY:
			err = handleDelEntity(client)
		case entitiesd.CMD_RPC:
			err = handleRPCToClient(client)
		}

		if err != nil {
			if IsNetworkError(err) {
				HandleConnectionError(client, err)
				break
			} else {
				log.Printf("Error while handling entitiesd: %s", err.Error())
			}
		}
	}
}

func handleNewEntity(client *entitiesd.EntitiesdClient) error {
	var clientid ClientId
	var eid Eid
	var entityType string

	client.RecvCid(&clientid)
	client.RecvEid(&eid)
	if err := client.RecvString(&entityType); err != nil {
		return err
	}
	log.Printf("handleNewEntity to %s, creating %s<%s>", clientid, entityType, eid)
	return dispatchNewEntityToClient(clientid, eid, entityType)
}

func handleDelEntity(client *entitiesd.EntitiesdClient) error {
	var clientid ClientId
	var eid Eid
	client.RecvCid(&clientid)
	if err := client.RecvEid(&eid); err != nil {
		return err
	}

	log.Printf("handleDelEntity to %s, deleting %s", clientid, eid)
	return dispatchDelEntityToClient(clientid, eid)
}

func handleRPCToClient(client *entitiesd.EntitiesdClient) error {
	var clientid ClientId
	var eid Eid
	var method string
	var args []interface{}

	client.RecvCid(&clientid)
	if err := client.RecvRPC(&eid, &method, &args); err != nil {
		return err
	}
	log.Printf("handleRPCToClient %s.%s, calling %s%v", clientid, eid, method, args)
	return dispatchRPCToClient(clientid, eid, method, args)
}
