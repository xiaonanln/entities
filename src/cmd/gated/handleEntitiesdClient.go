package main

import (
	. "common"
	"entitiesd"
	"log"
)

func handleNewEntity(client *entitiesd.EntitiesdClient) error {
	var clientid ClientId
	var eid Eid
	var entityType string
	client.RecvCid(&clientid)
	client.RecvEid(&eid)
	err := client.RecvString(&entityType)
	if err != nil {
		return err
	}

	log.Printf("handleNewEntity to %s, creating %s<%s>", clientid, entityType, eid)
	return nil
}

func handleDelEntity(client *entitiesd.EntitiesdClient) error {
	var clientid ClientId
	var eid Eid
	client.RecvCid(&clientid)
	if err := client.RecvEid(&eid); err != nil {
		return err
	}

	log.Printf("handleDelEntity to %s, deleting %s", clientid, eid)
	return nil
}
