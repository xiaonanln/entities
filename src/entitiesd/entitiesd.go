package entitiesd

import (
	. "common"
	"conf"
	"entities"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	pid int
)

func RegisterEntity(entity entities.EntityType) {
	entities.RegisterEntity(entity)
}

func NewEntity(entityType string) *entities.Entity {
	return entities.NewEntity(entityType)
}

func Run() {
	log.Println("Starting entitiesd service ...")
	parseArguments()
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
		switch cmd {
		case CMD_NEW_CLIENT:
			err = handleNewClient(client)
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
	var cid ClientId
	err := gated.RecvCid(&cid)
	if err != nil {
		return err
	}

	log.Printf("%s: cid %s", gated, cid)
	boot := newBootEntity()

	client := NewClient(gated.Gid, cid)
	boot.SetClient(client)
	return nil
}

func handleRPC(gated *EntitiesdClientProxy) error {
	var eid Eid
	var method string
	var args []interface{}
	err := gated.RecvRPC(&eid, &method, &args)
	if err != nil {
		return err
	}

	// received rpc from gate
	return nil
}

func newBootEntity() *entities.Entity {
	config := conf.GetEntitiesConfig()
	return NewEntity(config.BootEntity)
}
