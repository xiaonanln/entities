package entities

import (
	. "common"
	"fmt"
	"log"
	"reflect"
)

type callQueueItem struct {
	src    Eid // source entity
	method string
	args   []interface{}
}

func (self callQueueItem) String() string {
	return fmt.Sprintf(`CALL<%s(%v) FROM %s>`, self.method, self.args, self.src)
}

type EntityType interface {
	Id() Eid
	Call(id Eid, method string, args ...interface{})
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
func (self *Entity) Call(id Eid, method string, args ...interface{}) {
	entity := getEntity(id)
	if entity != nil {
		entity.pushCall(self.id, method, args)
		return
	}
	// call the coordinator now...
	log.Printf("entity %s not found, using coordinator...", id)
}

func (self *Entity) pushCall(srcId Eid, method string, args []interface{}) {
	self.callQueue <- &callQueueItem{src: srcId, method: method, args: args} // push call to call queue
}

func (self *Entity) SetClient(client *Client) {
	if self.Client != nil {
		self.Client.Close()
	}

	self.Client = client
	if self.Client != nil {
		self.Client.NewEntity(self)
	}
}

func (self *Entity) routine() {
	for {
		call := <-self.callQueue
		log.Printf("%s >>> %s.%s%v", call.src, self.id, call.method, call.args)
		self.handleCall(call)
	}
}

func (self *Entity) handleCall(call *callQueueItem) {
	defer func() {
		if err := recover(); err != nil {
			// error recovered in handleCall
			log.Printf("%s failed: %s", call, err)
		}
	}()
	method := self.realEntity.MethodByName(call.method)
	in := make([]reflect.Value, len(call.args))

	for i, arg := range call.args {
		in[i] = reflect.ValueOf(arg)
	}
	method.Call(in)
}
