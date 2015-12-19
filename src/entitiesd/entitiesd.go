package entitiesd

import (
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
		panic(err)
	}

	log.Println("Serving entitiesd service on", addr, "...")
	for {
		conn, err := servsock.Accept()
		if err != nil {
			panic(err)
		}
		log.Println("New connection:", conn.RemoteAddr())
		go serveClientConnection(conn)
	}
}

func serveClientConnection(conn net.Conn) {
	conn.Close()
}
