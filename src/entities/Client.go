package entities

import (
	. "common"
	"log"
)

type ClientRPCer interface {
	Call(id Eid, method string, args ...interface{})
	NewEntity(id Eid, entityType string)
	DelEntity(id Eid)
	Close()
}

type Client struct {
	rpcer ClientRPCer
}

func NewClient(rpcer ClientRPCer) *Client {
	return &Client{
		rpcer: rpcer,
	}
}

func (self *Client) Close() {
	self.rpcer.Close()
}

func (self *Client) Call(eid Eid, method string, args ...interface{}) {
	self.rpcer.Call(eid, method, args...)
}

func (self *Client) NewEntity(entity *Entity) {
	log.Println("NewEntity", entity, entity.EntityType())
	self.rpcer.NewEntity(entity.id, entity.EntityType())
}

func (self *Client) DelEntity(eid Eid) {
	self.rpcer.DelEntity(eid)
}
