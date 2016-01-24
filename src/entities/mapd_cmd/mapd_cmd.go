package mapd_cmd

import (
	. "common"
	"conf"
	"log"
	"mapd"
	"setup"
	"sync"
)

var (
	pid           int
	clientLock    sync.Mutex
	mapdCmdClient *mapd.MapdClient
)

func Init(_pid int) {
	if !setup.IsEntitiesd() {
		log.Fatalf("mapd_cmd.Init should only be called by entitiesd")
	}

	pid = _pid

	// clientLock.Lock()
	// defer clientLock.Unlock()
	// getMapdCmdClient() // connect
}

func getMapdCmdClient() (*mapd.MapdClient, error) {
	if mapdCmdClient == nil {
		if err := tryConnectMapdCmdService(); err != nil {
			return nil, err
		}
	}

	return mapdCmdClient, nil
}

func tryConnectMapdCmdService() error {
	mapdConfig := conf.GetMapdConfig()
	log.Printf("Connecting mapd.cmd @%v:%v ...", mapdConfig.Host, mapdConfig.CmdPort)
	conn, err := ConnectTCP(mapdConfig.Host, mapdConfig.CmdPort)
	if err != nil {
		return err
	}
	mapdCmdClient = mapd.NewMapdClient(conn)

	if err = mapdCmdClient.SetPid(pid); err != nil {
		onMapdCmdClientError(err)
		return err
	}

	return nil
}

func DeclareNewEntity(eid Eid) error {
	clientLock.Lock()
	defer clientLock.Unlock()

	client, err := getMapdCmdClient()
	if err != nil {
		return err
	}

	if err = client.SetMapping(eid); err != nil {
		onMapdCmdClientError(err)
		return err
	}

	return nil
}

func RegisterGlobalEntity(eid Eid, entityType string) (bool, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	client, err := getMapdCmdClient()
	if err != nil {
		return false, err
	}
	ok, err := client.RegisterGlobal(eid, entityType)
	if err != nil {
		onMapdCmdClientError(err)
	}
	return ok, err
}

func RPC(eid Eid, method string, args []interface{}) error {
	clientLock.Lock()
	defer clientLock.Unlock()

	client, err := getMapdCmdClient()
	if err != nil {
		return err
	}
	err = client.RPC(eid, method, args)
	if err != nil {
		return err
	}
	return nil
}

func onMapdCmdClientError(err error) {
	HandleConnectionError(mapdCmdClient, err)
	mapdCmdClient.Close()
	mapdCmdClient = nil
}
