package main

import (
	. "common"
	"conf"
	"flag"
	"fmt"
	"gated"
	"log"
	"net"
	"os"
	"time"
)

var (
	accountId Eid
	avatarId  Eid
	gid       int
)

func main() {
	gatedCount := conf.GetGatedCount()
	log.Println("Found", gatedCount, "gated")
	for gid := 1; gid <= gatedCount; gid++ {
		gatedConfig := conf.GetGatedConfig(gid)
		log.Printf("Gated %d config: %v", gid, gatedConfig)
	}

	parseArguments()
	gatedConfig := conf.GetGatedConfig(gid)
	log.Printf("Connecting to gated %d, config %v", gid, gatedConfig)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", gatedConfig.Port))
	if err != nil {
		HandleConnectionError(conn, err)
		os.Exit(1)
	}

	log.Printf("Connect successfully")
	client := gated.NewGatedClient(conn)

	go sendRoutine(client)
	receiveRoutine(client)
	client.Close()
}

func receiveRoutine(client *gated.GatedClient) {
	for {
		cmd, err := client.RecvCmd()
		if err != nil {
			HandleConnectionError(client, err)
			break
		}

		log.Printf("Received cmd: %v", cmd)

		switch cmd {
		case gated.CMD_NEW_ENTITY:
			handleNewEntity(client)
		case gated.CMD_DEL_ENTITY:
			handleDelEntity(client)
		case gated.CMD_RPC:
			handleRPC(client)
		default:
			log.Printf("unknown cmd: %v", cmd)
			client.Close()
			break
		}
		// var eid Eid
		// var method string
		// var args []interface{}
		// err = client.RecvRPC(&eid, &method, &args)
		// if err != nil {
		// 	HandleConnectionError(client, err)
		// 	break
		// }

		// log.Printf("RECV RPC: %s.%s%v", eid, method, args)
	}
}

func handleNewEntity(client *gated.GatedClient) error {
	var eid Eid
	var entityType string
	client.RecvEid(&eid)
	if err := client.RecvString(&entityType); err != nil {
		return err
	}
	log.Printf("NewEntity: %s, entityType %s", eid, entityType)
	if entityType == "Account" {
		accountId = eid
	} else if entityType == "Avatar" {
		avatarId = eid
	} else {
		log.Printf("invalid entity type: %s", entityType)
	}
	return nil
}

func handleDelEntity(client *gated.GatedClient) error {
	var eid Eid
	if err := client.RecvEid(&eid); err != nil {
		return err
	}
	log.Printf("DelEntity: %s", eid)
	if accountId == eid {
		accountId = ""
		log.Println("Account destroyed")
	} else if avatarId == eid {
		avatarId = ""
		log.Println("Avatar destroyed")
	}
	return nil
}

func handleRPC(client *gated.GatedClient) error {
	var eid Eid
	var method string
	var args []interface{}
	if err := client.RecvRPC(&eid, &method, &args); err != nil {
		return err
	}
	log.Printf("RecvRPC: %s.%s%v", eid, method, args)
	return nil
}

func sendRoutine(client *gated.GatedClient) {
	const (
		QUERY_ROOMS = iota
	)

	nextAction := QUERY_ROOMS

	for {
		time.Sleep(time.Second * 3)

		if avatarId == "" && accountId != "" {
			client.Call(accountId, "Login", "test", "test")
			continue
		}

		if avatarId == "" {
			log.Println("Waiting for login to complete ...")
			continue
		}

		if nextAction == QUERY_ROOMS {
			client.Call(avatarId, "QueryChatRoomList")
			continue
		}

		log.Println("Found no account and no avatar, do't know what to do ...")
	}
}

func parseArguments() {
	flag.IntVar(&gid, "gid", 1, "Connect specified gated. Default to be 1")
	flag.Parse()
}
