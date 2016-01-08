package main

import (
	. "common"
	"conf"
	"flag"
	"fmt"
	"gated"
	"log"
	"net"
	_ "setup"
	"sync"
)

var (
	gid          int
	gatedConfig  *conf.GatedConfig
	waitServices = sync.WaitGroup{}
)

func init() {
	log.Println("gated main init...")
}

func main() {
	parseArguments()
	readConfig()
	startGatedServices()
	waitServices.Wait()
}

func parseArguments() {
	flag.IntVar(&gid, "gid", 0, "process id, should be unique for entitiesd services")
	flag.Parse()

	log.Println("Gid:", gid)
	if gid <= 0 {
		panic("gid must be positive")
	}
}

func readConfig() {
	gatedConfig = conf.GetGatedConfig(gid)
	log.Printf("Using gated config: %v\n", *gatedConfig)
}

func startGatedServices() {
	waitServices.Add(1) // number of services
	go serveConnectionService()
}

func serveConnectionService() {
	addr := fmt.Sprintf("%s:%d", gatedConfig.Host, gatedConfig.Port)
	servsock, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Serving gated cmd service on", addr, "...")
	for {
		conn, err := servsock.Accept()
		if err != nil {
			panic(err)
		}
		log.Println("New connection:", conn.RemoteAddr())
		go serveClientConnection(conn)
	}
	waitServices.Done()
}

func serveClientConnection(conn net.Conn) {
	gatedClient := gated.NewGatedClientProxy(conn)
	defer gatedClient.Close()

	// put gated to dispatcher
	dispatchAddNewClient(gatedClient)
	defer dispatchOnClientClose(gatedClient)

	err := onClientConnect(gatedClient)
	if err != nil {
		HandleConnectionError(gatedClient, err)
		return
	}

	serveClientConnectionLoop(gatedClient)
}

func serveClientConnectionLoop(gatedClient *gated.GatedClientProxy) {
	for {
		cmd, err := gatedClient.RecvCmd()
		if err != nil {
			HandleConnectionError(gatedClient, err)
			break
		}

		switch cmd {
		case gated.CMD_RPC:
			var eid Eid
			var method string
			var args []interface{}
			err = gatedClient.RecvRPC(&eid, &method, &args)
			if err != nil {
				HandleConnectionError(gatedClient, err)
				break
			}
			err = onClientCallRPC(gatedClient, eid, method, args)
			if err != nil {
				HandleConnectionError(gatedClient, err)
				break
			}
		default:
			log.Println("Invalid cmd: %s", cmd)
			break
		}
	}
}
