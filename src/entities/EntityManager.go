package entities

import (
	. "common"
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
)

func NewEntity(entityTypeName string) *Entity {
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

func OnCall(clientid ClientId, eid Eid, method string, args []interface{}) {
	entity := getEntity(eid)
	if entity == nil {
		log.Printf("entity %s not found when calling %s%v", eid, method, args)
		return
	}

	entity.pushCall(eid, method, args) // calling from client is calling from self
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
