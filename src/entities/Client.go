package entities

import (
	. "common"
	"conf"
	"log"
	"sync"
)

var (
	clientsLock sync.RWMutex
	clients     = make(map[ClientId]*Client)
)

type ClientRPCer interface {
	Call(id Eid, method string, args ...interface{})
	NewEntity(id Eid, entityType string)
	DelEntity(id Eid)
	Close()
}

type Client struct {
	clientid ClientId
	rpcer    ClientRPCer
	owner    Eid
}

func newClient(clientid ClientId, rpcer ClientRPCer) *Client {
	client := &Client{
		clientid: clientid,
		rpcer:    rpcer,
	}
	clientsLock.Lock()
	defer clientsLock.Unlock()
	clients[clientid] = client
	return client
}

func getClient(clientid ClientId) *Client {
	clientsLock.RLock()
	defer clientsLock.RUnlock()
	return clients[clientid]
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

func (self *Client) setOwner(owner Eid) {
	self.owner = owner
}

func newBootEntity() *Entity {
	config := conf.GetEntitiesConfig()
	return NewEntity(config.BootEntity)
}

func OnNewClient(clientid ClientId, rpcer ClientRPCer) error {
	boot := newBootEntity()

	client := newClient(clientid, rpcer)
	boot.SetClient(client)
	return nil
}

func OnDelClient(clientid ClientId) error {
	client := getClient(clientid)
	if client == nil {
		log.Printf("OnDelClient: client %s not found", clientid)
		return nil
	}
	if client.owner != "" {
		if entity := getEntity(client.owner); entity != nil {
			entity.SetClient(nil)
		}
	}

	return nil
}
