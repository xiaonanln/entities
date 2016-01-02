package main

import (
	. "common"
	"conf"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	gid int
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

}

func parseArguments() {
	flag.IntVar(&gid, "gid", 1, "Connect specified gated. Default to be 1")
	flag.Parse()
}
