package main

import (
	. "common"
	"log"
	"mapd"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connect success")
	client := mapd.NewMapdClient(conn)
	log.Println("Sending pid 1...")
	err = client.SetPid(1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Set pid successfully")

	eid := NewEid()
	err = client.SetMapping(eid)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Set mapping successfully")

	client.RPC(eid, "TestRPC", []interface{}{1, "2", 3.0})
	log.Println("Send RPC successfully")
	client.Close()
}
