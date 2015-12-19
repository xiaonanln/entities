package main

import (
	"entitiesd"
	"sample_entities"
)

func main() {
	entitiesd.RegisterEntity(&entitiesd.TestEntity{})
	entitiesd.RegisterEntity(&sample_entities.Account{})
	entitiesd.Run()
}
