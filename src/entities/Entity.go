package entities

import (
	. "common"

	"entities/mapd_cmd"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
)

var (
	interfaceSliceType     = reflect.TypeOf([]interface{}{})          // get the type []interface{}
	stringInterfaceMapType = reflect.TypeOf(map[string]interface{}{}) // get the type map[string]interface{}
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

type ClientEventHandler interface {
	OnLoseClient(oldClient *Client)
	OnGetClient()
}

type Entity struct {
	id                 Eid
	callQueue          chan *callQueueItem
	realEntity         EntityType
	realEntityValue    reflect.Value
	clientEventHandler ClientEventHandler
	Client             *Client
}

func (self *Entity) init(id Eid, realEntityValue reflect.Value) {
	self.id = id
	self.realEntityValue = realEntityValue
	self.realEntity = realEntityValue.Interface().(EntityType)
	self.clientEventHandler, _ = self.realEntity.(ClientEventHandler)
	self.callQueue = make(chan *callQueueItem, ENTITY_CALL_QUEUE_BUFF_LEN)
}

func (self *Entity) Id() Eid {
	return self.id
}

func (self *Entity) EntityType() string {
	return reflect.Indirect(self.realEntityValue).Type().Name()
}

func (self *Entity) String() string {
	return fmt.Sprintf(`%s<%s>`, self.EntityType(), self.Id())
}

func (self *Entity) OnInit() {
	log.Printf("WARNING: %s.OnInit is not defined", self)
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

func (self *Entity) CallGlobalEntity(entityType string, method string, args ...interface{}) error {
	eid, ok := globalRegisterMap[entityType]

	if !ok {
		log.Printf("ERROR: Global entity %s not found while calling method %s", entityType, method)
		return fmt.Errorf("global entity %s not found", entityType)
	}

	return self.Call(eid, method, args...)
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
	if self.clientEventHandler != nil {
		self.clientEventHandler.OnGetClient()
	}
}

func (self *Entity) onLoseClient(oldClient *Client) {
	log.Printf("Entity %s lose client %s", self, oldClient)
	if self.clientEventHandler != nil {
		self.clientEventHandler.OnLoseClient(oldClient)
	}
}

func (self *Entity) Destroy() {
	if self.Client != nil {
		self.SetClient(nil)
	}
	close(self.callQueue)
	delEntity(self.id)
}

func (self *Entity) NewEntity(entityType string) (*Entity, error) {
	return NewEntity(entityType)
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
	method := self.realEntityValue.MethodByName(call.method)
	methodType := method.Type()

	in := make([]reflect.Value, len(call.args))

	for i, arg := range call.args {
		argType := methodType.In(i)
		argVal := reflect.ValueOf(arg)
		in[i] = self.convertMethodArgType(argVal, argType)
		// log.Printf("Arg %d is %T %v value %v => %v", i, arg, arg, argVal, in[i])
	}
	// log.Printf("arguments: %v", in)
	method.Call(in)
}

func (self *Entity) convertMethodArgType(val reflect.Value, typ reflect.Type) reflect.Value {
	valType := val.Type()
	if valType.ConvertibleTo(typ) {
		return val.Convert(typ)
	}
	// can not convert directly
	intSlice, ok := val.Interface().([]interface{})
	if ok {
		// val is of type []interface{}, try to convert to typ
		targetVal := reflect.New(typ)

	}
}
