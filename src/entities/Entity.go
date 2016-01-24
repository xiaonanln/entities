package entities

import (
	. "common"

	"entities/mapd_cmd"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
)

type callQueueItem struct {
	method string
	args   []interface{}
}

func (self callQueueItem) String() string {
	return fmt.Sprintf(`CALL<%s(%v)>`, self.method, self.args)
}

type EntityType interface {
	Id() Eid
	EntityType() string
	Call(id Eid, method string, args ...interface{}) error
	SetClient(client *Client)
	GiveClientTo(other *Entity)
}

type Entity struct {
	id         Eid
	callQueue  chan *callQueueItem
	realEntity reflect.Value
	Client     *Client
}

func (self *Entity) init(id Eid, realEntity reflect.Value) {
	self.id = id
	self.realEntity = realEntity
	self.callQueue = make(chan *callQueueItem, ENTITY_CALL_QUEUE_BUFF_LEN)
}

func (self *Entity) Id() Eid {
	return self.id
}

func (self *Entity) EntityType() string {
	return reflect.Indirect(self.realEntity).Type().Name()
}

func (self *Entity) String() string {
	return fmt.Sprintf(`%s<%s>`, self.EntityType(), self.Id())
}

// call another entity
func (self *Entity) Call(id Eid, method string, args ...interface{}) error {
	// entity := getEntity(id)
	// if entity != nil {
	// 	entity.pushCall(method, args)
	// 	return nil
	// }

	// call the coordinator now...
	log.Printf("entity %s not found, using mapd...", id)
	err := mapd_cmd.RPC(id, method, args)
	return err
}

func (self *Entity) pushCall(method string, args []interface{}) {
	self.callQueue <- &callQueueItem{method: method, args: args} // push call to call queue
}

func (self *Entity) SetClient(client *Client) {
	oldClient := self.Client

	if self.Client != nil {
		self.Client.owner = ""
		self.Client.Close()
	}

	self.Client = client
	if client != nil {
		client.owner = self.id
		client.NewEntity(self)
	}

	if oldClient != nil && client == nil {
		self.onLoseClient(oldClient)
	} else if oldClient == nil && client != nil {
		self.onGetClient()
	}
}

func (self *Entity) GiveClientTo(other *Entity) {
	if self.Client == nil || self == other {
		return
	}

	client := self.Client
	client.owner = ""
	self.Client = nil

	client.DelEntity(self.id)
	self.onLoseClient(client)

	other.SetClient(client)
	return
}

func (self *Entity) onGetClient() {
	log.Printf("Entity %s get client %s", self, self.Client)
}

func (self *Entity) onLoseClient(oldClient *Client) {
	log.Printf("Entity %s lose client %s", self, oldClient)
}

func (self *Entity) Destroy() {
	if self.Client != nil {
		self.SetClient(nil)
	}
	delEntity(self.id)
}

func (self *Entity) routine() {
	for {
		call := <-self.callQueue
		log.Printf(">>> %s.%s%v", self.id, call.method, call.args)
		self.handleCall(call)
	}
}

func (self *Entity) handleCall(call *callQueueItem) {
	defer func() {
		if err := recover(); err != nil {
			// error recovered in handleCall
			log.Printf("%s failed: %s", call, err)
			debug.PrintStack()
		}
	}()
	method := self.realEntity.MethodByName(call.method)
	in := make([]reflect.Value, len(call.args))

	for i, arg := range call.args {
		in[i] = reflect.ValueOf(arg)
	}
	method.Call(in)
}
