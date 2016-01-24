package mapd_rpc

import (
	. "common"
	"conf"
	"entities"
	"log"
	"mapd"
	"setup"
	"time"
)

var (
	pid           int
	mapdRpcClient *mapd.MapdClient
)

func Init(_pid int) {
	if !setup.IsEntitiesd() {
		log.Fatalf("mapd_cmd.Init should only be called by entitiesd")
	}

	pid = _pid
	go maintainMapdRpcClient()
}

func maintainMapdRpcClient() {
	var err error
	for {
		if mapdRpcClient == nil {
			err = tryConnectMapdRpcService()
			if err != nil {
				log.Printf("Connect mapd.RPC failed: %s", err)
				time.Sleep(time.Second)
				continue
			}
		}

		var eid Eid
		var method string
		var args []interface{}
		err = mapdRpcClient.RecvRPC(&eid, &method, &args)
		if err != nil {
			onMapdRpcClientError(err)
			continue
		}

		entities.OnCall(eid, method, args)
	}
}

func tryConnectMapdRpcService() error {
	mapdConfig := conf.GetMapdConfig()
	log.Printf("Connecting mapd.RPC @%v:%v ...", mapdConfig.Host, mapdConfig.RPCPort)
	conn, err := ConnectTCP(mapdConfig.Host, mapdConfig.RPCPort)
	if err != nil {
		return err
	}
	mapdRpcClient = mapd.NewMapdClient(conn)

	if err = mapdRpcClient.SendPid(pid); err != nil {
		onMapdRpcClientError(err)
		return err
	}

	return nil
}

func onMapdRpcClientError(err error) {
	if mapdRpcClient != nil {
		HandleConnectionError(mapdRpcClient, err)
	}
	mapdRpcClient.Close()
	mapdRpcClient = nil
}
