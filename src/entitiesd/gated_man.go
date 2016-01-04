package entitiesd

import (
	. "common"
	"conf"
	"log"
	"setup"
	"sync/atomic"
	"unsafe"
)

type gatedSendOp struct {
	op  func(gated *EntitiesdClientProxy) error
	gid int
}

var (
	gateds          []*EntitiesdClientProxy
	gatedSendOpChan = make(chan gatedSendOp)
)

func init() {
	if setup.IsEntitiesd() {
		gateds = make([]*EntitiesdClientProxy, conf.GetGatedCount())
		log.Printf("gated_man: found %d gated", len(gateds))
		go serviceGatedSendOperations()
	}
}

func serviceGatedSendOperations() {
	for {
		sendOp := <-gatedSendOpChan
		gid := sendOp.gid
		gated := getGated(gid)
		if gated == nil {
			// what we do when gated is down ?
			log.Printf("Gated %d is not online, RPCs have to be dropped!", gid)
			continue
		}

		err := sendOp.op(gated)
		if err != nil {
			HandleConnectionError(gated, err)
			log.Printf("Gated %d send error: %s, RPCs have to be dropped", gid, err)
		}
	}
}

func onNewGated(gated *EntitiesdClientProxy, gid int) {
	oldGated := getGated(gid)
	if oldGated != nil {
		oldGated.Close()
	}

	gated.SetGid(gid)
	setGated(gid, gated)
	log.Printf("New gated %s for gid %d", gated, gid)
}

func getGated(gid int) *EntitiesdClientProxy {
	addr := unsafe.Pointer(&gateds[gid-1])
	ptr := atomic.LoadPointer(&addr)
	return (*EntitiesdClientProxy)(ptr)
}

func setGated(gid int, gated *EntitiesdClientProxy) {
	addr := unsafe.Pointer(&gateds[gid-1])
	atomic.StorePointer(&addr, unsafe.Pointer(gated))
}

func onCallClient(gid int, clientid ClientId, eid Eid, method string, args []interface{}) {
	gatedSendOpChan <- gatedSendOp{
		op: func(gated *EntitiesdClientProxy) error {
			return gated.RPC(clientid, eid, method, args)
		},
		gid: gid,
	}
}

func onNewEntity(gid int, clientid ClientId, eid Eid, entityType string) {
	gatedSendOpChan <- gatedSendOp{
		op: func(gated *EntitiesdClientProxy) error {
			return gated.NewEntity(clientid, eid, entityType)
		},
		gid: gid,
	}
}

func onDelEntity(gid int, clientid ClientId, eid Eid) {
	gatedSendOpChan <- gatedSendOp{
		op: func(gated *EntitiesdClientProxy) error {
			return gated.DelEntity(clientid, eid)
		},
		gid: gid,
	}
}
