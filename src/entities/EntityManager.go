package entities

import (
	. "common"
	"entities/mapd_cmd"
	"fmt"
	"log"
	"reflect"
	"sync"
)

const (
	ENTITY_CALL_QUEUE_BUFF_LEN = 0
)

var (
	entitiesLock          sync.RWMutex
	entities              = make(map[Eid]*Entity)
	registeredEntityTypes = make(map[string]reflect.Type)
	globalRegisterMap     = make(map[string]Eid)
)

func newEntity(entityTypeName string) *Entity {
	entityType, ok := registeredEntityTypes[entityTypeName]
	if !ok {
		panic(fmt.Errorf("unknown entity type: %s", entityTypeName))
	}

	entityPtr := reflect.New(entityType)
	entityValue := reflect.Indirect(entityPtr)
	entity := entityValue.FieldByName("Entity").Addr().Interface().(*Entity)
	// setup inner entity
	id := NewEid()
	entity.init(id, entityPtr)

	entities[id] = entity
	go entity.routine()
	return entity
}

func NewEntity(entityType string) (*Entity, error) {
	entity := newEntity(entityType)
	if entity == nil {
		return nil, fmt.Errorf("NewEntity failed, type: %s", entityType)
	}

	entity.pushCall("OnInit", []interface{}{})

	err := mapd_cmd.DeclareNewEntity(entity.id)
	if err != nil {
		// declare new entity failed, we have to
		log.Printf("Declare new entity %s failed: %s", entity, err)
		entity.Destroy()
		return nil, err
	}

	return entity, nil
}

func NewGlobalEntity(entityType string) (*Entity, error) {
	entity, err := NewEntity(entityType)
	if err != nil {
		return entity, err
	}

	ok, err := mapd_cmd.RegisterGlobalEntity(entity.id, entityType)
	if err != nil {
		entity.Destroy()
		return nil, err
	}

	if ok {
		OnGlobalRegister(entityType, entity.id)
		entity.pushCall("OnRegisteredGlobally", []interface{}{})
		return entity, nil // register global succeed
	} else {
		entity.Destroy() // register global failed, has to destroy entity
		return nil, nil  // no error and no entity
	}
}

func getEntity(id Eid) *Entity {
	entitiesLock.RLock()
	defer entitiesLock.RUnlock()
	return entities[id]
}

func putEntity(entity *Entity) {
	entitiesLock.Lock()
	defer entitiesLock.Unlock()
	if _, ok := entities[entity.id]; ok {
		log.Panicf("entity %s already exists", entity.id)
	}

	entities[entity.id] = entity
}

func delEntity(eid Eid) {
	entitiesLock.Lock()
	defer entitiesLock.Unlock()
	delete(entities, eid)
}

func RegisterEntity(entity EntityType) {
	entityValue := reflect.Indirect(reflect.ValueOf(entity))
	entityType := entityValue.Type()
	typeName := entityType.Name()
	entityField := entityValue.FieldByName("Entity") // User entity class must have Entity field defined as Entity Type
	if !entityField.IsValid() {
		panic(fmt.Errorf("%s is not a valid Entity type", typeName))
	}

	_, exists := registeredEntityTypes[typeName]
	if exists {
		panic(fmt.Errorf("entity type %s is aleady registered", typeName))
	}

	registeredEntityTypes[typeName] = entityType
	log.Printf("entity type %s registered successfully: %v", typeName, entityType)
	// fmt.Println(entityValue, entityField, entityField.IsValid())
}

func OnCall(eid Eid, method string, args []interface{}) {
	entity := getEntity(eid)
	if entity == nil {
		log.Printf("entity %s not found when calling %s%v", eid, method, args)
		return
	}

	entity.pushCall(method, args) // calling from client is calling from self
}

func OnGlobalRegister(entityType string, eid Eid) {
	globalRegisterMap[entityType] = eid
}

func getLocalGlobalEntity(entityType string) *Entity {
	eid, ok := globalRegisterMap[entityType]
	if !ok {
		log.Printf("WARNING: global entity %s is not registered", entityType)
		return nil
	}
	return getEntity(eid)
}

// func setLocalEntity(ent *Entity) {
// 	id := ent.id
// 	_, ok := entities[id]
// 	if ok {
// 		// entity already exists with same id
// 		log.Printf("entity %s already exists in EntityManager\n", id)
// 		return
// 	}
// 	entities[id] = ent
// }
