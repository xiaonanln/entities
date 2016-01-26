package main

import (
	. "common"
	"mapd"
	"sync"
)

var (
	registerGlobalLock       sync.RWMutex
	registeredGlobalEntities = make(map[string]Eid)
)

func registerGlobalEntity(eid Eid, entityType string) bool {
	registerGlobalLock.Lock()
	defer registerGlobalLock.Unlock()
	if _, ok := registeredGlobalEntities[entityType]; ok {
		// already registered
		return false
	}

	registeredGlobalEntities[entityType] = eid
	return true
}

func getRegisteredGlobalEntity(entityType string) Eid {
	registerGlobalLock.RLock()
	defer registerGlobalLock.RUnlock()
	return registeredGlobalEntities[entityType]
}

func NotifyAllRegisteredGlobalEntities(client *mapd.MapdClientProxy) {
	lock.RLock()
	defer lock.RUnlock()

	for entityType, eid := range registeredGlobalEntities {
		client.NotifyRegisterGlobal(eid, entityType)
	}
}
