package mapd_man

import (
	. "common"
	"conf"
	"entities"
	"mapd"
	"setup"
	"sync"
)

var (
	accessLock    sync.Mutex
	mapdCmdClient *mapd.MapdClient
	mapdRpcClient *mapd.MapdClient
)

func init() {
	if setup.IsEntitiesd() {
		go maintainMapdClient()
	}
}

func maintainMapdClient() {

}

func getMapdCmdClient() {
	if mapdCmdClient == nil {
		tryConnectMapdCmdService()
	}
}

func tryConnectMapdCmdService() error {
	mapdConfig := conf.GetMapdConfig()
	conn, err := ConnectTCP(mapdConfig.Host, mapdConfig.CmdPort)
	if err != nil {
		return err
	}
	mapdCmdClient = mapd.NewMapdClient(conn)
	return nil
}

func tryConnectMapdRpcService() error {
	mapdConfig := conf.GetMapdConfig()
	conn, err := ConnectTCP(mapdConfig.Host, mapdConfig.RPCPort)
	if err != nil {
		return err
	}
	mapdRpcClient = mapd.NewMapdClient(conn)
	return nil
}

func SetNewEntity(entity *entities.Entity) {

}
