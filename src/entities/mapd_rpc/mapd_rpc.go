package mapd_rpc

import (
	. "common"
	"conf"
	"log"
	"mapd"
	"setup"
	"time"
)

var (
	onCall               func(eid Eid, method string, args []interface{})
	pid                  int
	mapdRpcClient        *mapd.MapdClient
	globalEntityRegister = make(map[string]Eid)
)

func Init(_pid int, _onCall func(eid Eid, method string, args []interface{})) {
	if !setup.IsEntitiesd() {
		log.Fatalf("mapd_cmd.Init should only be called by entitiesd")
	}

	pid = _pid
	onCall = _onCall
	go maintainMapdRpcClient()
}

func GetRegisteredGlobalEntity(entityType string) Eid {
	return globalEntityRegister[entityType]
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

		cmd, err := mapdRpcClient.RecvCmd()
		log.Printf("From mapd <<< cmd %v error %s", cmd, err)
		switch cmd {
		case mapd.CMD_RPC:

			var eid Eid
			var method string
			var args []interface{}
			err = mapdRpcClient.RecvRPC(&eid, &method, &args)
			if err != nil {
				onMapdRpcClientError(err)
				continue
			}

			onCall(eid, method, args)

		case mapd.CMD_REGISTER_GLOBAL:
			var eid Eid
			var entityType string
			mapdRpcClient.RecvEid(&eid)
			err := mapdRpcClient.RecvString(&entityType)
			if err != nil {
				onMapdRpcClientError(err)
				continue
			}

			oldEid, ok := globalEntityRegister[entityType]
			if ok {
				// global entity type already registered
				log.Printf("ERROR: Global entity type %s is already registered to entity %s !!!!!!", entityType, oldEid)
			}

			globalEntityRegister[entityType] = eid
			log.Printf("Global entity %s registered to be %s", entityType, eid)
		}
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
