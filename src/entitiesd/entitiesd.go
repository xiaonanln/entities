package entitiesd

import (
	. "common"
	"conf"
	"entities"
	"entities/mapd_cmd"
	"entities/mapd_rpc"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	pid             int
	NewEntity       = entities.NewEntity
	NewGlobalEntity = entities.NewGlobalEntity
)

func RegisterEntity(entity entities.EntityType) {
	entities.RegisterEntity(entity)
}

func Prepare() {
	log.Println("Starting entitiesd service ...")
	parseArguments()
	log.SetPrefix(fmt.Sprintf("entitiesd-%d ", pid))
	mapd_cmd.Init(pid)
	mapd_rpc.Init(pid, entities.OnCall, entities.OnGlobalRegister)
}

func Run() {
	runEntitiesdService()
}

func parseArguments() {
	flag.IntVar(&pid, "pid", 0, "process id, should be unique for entitiesd services")
	flag.Parse()

	log.Println("Pid:", pid)
	if pid <= 0 {
		panic("pid must be positive")
	}
}

func runEntitiesdService() {
	config := conf.GetEntitiesdConfig(pid)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	servsock, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Listen error: %s", err)
	}

	log.Println("Serving entitiesd service on", addr, "...")
	for {
		conn, err := servsock.Accept()
		if err != nil {
			log.Fatalf("Accept error: %s", err)
		}

		log.Println("New connection:", conn.RemoteAddr())
		go serveClientConnection(conn)
	}
}

func serveClientConnection(conn net.Conn) {
	client := NewEntitiesdClientProxy(conn)
	defer client.Close()

	gid, err := client.RecvGid()
	if err != nil {
		HandleConnectionError(client, err)
		return
	}

	onNewGated(client, gid)

	serveClientConnectionLoop(client)
}

func serveClientConnectionLoop(client *EntitiesdClientProxy) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		// recovered from error, restart service loop
	// 		go serveClientConnectionLoop(client)
	// 	}
	// }()
	for {
		cmd, err := client.RecvCmd()
		if err != nil {
			HandleConnectionError(client, err)
			break
		}
		log.Printf("%s >>> cmd %v", client, cmd)
		switch cmd {
		case CMD_NEW_CLIENT:
			err = handleNewClient(client)
		case CMD_DEL_CLIENT:
			err = handleDelClient(client)
		case CMD_RPC:
			err = handleRPC(client)
		default:
			err = fmt.Errorf("Invalid cmd: %s", cmd)
		}

		if err != nil {
			HandleConnectionError(client, err)
			break
		}
	}
}

func handleNewClient(gated *EntitiesdClientProxy) error {
	var clientid ClientId
	err := gated.RecvCid(&clientid)
	if err != nil {
		return err
	}

	log.Printf("%s: new client %s", gated, clientid)
	clientRpcer := NewClientRPCProxy(gated.Gid, clientid)
	entities.OnNewClient(clientid, clientRpcer)

	return nil
}

func handleDelClient(gated *EntitiesdClientProxy) error {
	var clientid ClientId
	if err := gated.RecvCid(&clientid); err != nil {
		return err
	}
	log.Printf("%s: lose client %s", gated, clientid)
	entities.OnDelClient(clientid)
	return nil
}

func handleRPC(gated *EntitiesdClientProxy) error {
	var clientid ClientId
	var eid Eid
	var method string
	var args []interface{}

	gated.RecvCid(&clientid)

	err := gated.RecvRPC(&eid, &method, &args)
	if err != nil {
		return err
	}

	// received rpc from gate
	log.Printf("Client %s >>> %s.%s%v", clientid, eid, method, args)
	entities.OnCall(eid, method, args)

	return nil
}
