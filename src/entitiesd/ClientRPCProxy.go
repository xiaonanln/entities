package entitiesd

import (
	. "common"
	"log"
)

type ClientRPCProxy struct {
	Gid      int
	ClientId ClientId
}

func NewClientRPCProxy(gid int, clientid ClientId) *ClientRPCProxy {
	return &ClientRPCProxy{
		Gid:      gid,
		ClientId: clientid,
	}
}

func (self *ClientRPCProxy) Call(eid Eid, method string, args ...interface{}) {
	log.Printf("Call %s.%s%v", eid, method, args)
	onCallClient(self.Gid, self.ClientId, eid, method, args)
}

func (self *ClientRPCProxy) NewEntity(id Eid, entityType string) {
	log.Printf("NewEntity %s type %s", id, entityType)
	onNewEntity(self.Gid, self.ClientId, id, entityType)
}

func (self *ClientRPCProxy) DelEntity(id Eid) {
	log.Printf("DelEntity %s", id)
	onDelEntity(self.Gid, self.ClientId, id)
}

func (self *ClientRPCProxy) Close() {
}
