package entities

import (
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
	Id() string
}

type Entity struct {
	id         Eid
	callQueue  chan *callQueueItem
	realEntity reflect.Value
}

func (self *Entity) init(id Eid, realEntity reflect.Value) {
	self.id = id
	self.realEntity = realEntity
	self.callQueue = make(chan *callQueueItem, ENTITY_CALL_QUEUE_BUFF_LEN)
}

func (self *Entity) Id() Eid {
	return self.id
}

func (self *Entity) String() string {
	return fmt.Sprintf(`Entity<%s>`, self.id)
}

func (self *Entity) Call(id Eid, method string, args ...interface{}) {
	entity := GetLocalEntity(id)
	if entity != nil {
		entity.callQueue <- &callQueueItem{src: self.id, method: method, args: args} // push call to call queue
		return
	}
	// call the coordinator now...
	log.Printf("entity %s not found, using coordinator...", id)
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

type TestEntity struct {
	Entity
}

func (self *TestEntity) Foo(a int) {
	fmt.Println("RPC SUCCESSFULLY:", self, "Foo", a)
}

func (self *TestEntity) Bar(a string, b string) {
	fmt.Println("RPC:", self, "Bar", a, b)
}
