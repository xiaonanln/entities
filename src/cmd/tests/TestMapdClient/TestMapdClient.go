package main

import (
	"entities"
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

	eid := entities.NewEid()
	err = client.SetMapping(eid)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Set mapping successfully")

	client.Close()
}
