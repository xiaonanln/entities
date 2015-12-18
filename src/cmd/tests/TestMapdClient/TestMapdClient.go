package main

import (
	. "common"
	"log"
	"mapd"
	"net"
)

func main() {
	rpcConn, err := net.Dial("tcp", "localhost:5001")
	rpcClient := mapd.NewMapdClient(rpcConn)
	rpcClient.SendPid(1)

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

	var method string
	var args []interface{}
	err = rpcClient.RecvRPC(&eid, &method, &args)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Receive RPC successfully:", eid, method, args)

	client.Close()
	rpcClient.Close()
}
