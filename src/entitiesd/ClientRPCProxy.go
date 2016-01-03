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

func (self *ClientRPCProxy) Call(id Eid, method string, args ...interface{}) {

}

func (self *ClientRPCProxy) NewEntity(id Eid, entityType string) {
	log.Printf("NewEntity %s type %s", id, entityType)
}

func (self *ClientRPCProxy) DelEntity(id Eid) {
	log.Printf("DelEntity %s", id)
}

func (self *ClientRPCProxy) Close() {
}
