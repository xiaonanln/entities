package entities

import (
	"fmt"
	"log"
	"reflect"
)

const (
	ENTITY_CALL_QUEUE_BUFF_LEN = 0
)

var (
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

func GetLocalEntity(id Eid) *Entity {
	ent, ok := entities[id]
	if !ok {
		return nil
	}
	return ent
}

func RegisterEntity(entity interface{}) {
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
