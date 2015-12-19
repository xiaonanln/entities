package main

import (
	"entities"
	"entitiesd"
	"time"
)

func main() {
	entities.RegisterEntity(&entitiesd.TestEntity{})

	// entities.RegisterEntity(&entities.TestEntity{})
	e1 := entities.NewEntity("TestEntity")
	e2 := entities.NewEntity("TestEntity")
	e1.Call(e2.Id(), "Foo", 1)
	e1.Call(e2.Id(), "Bar", "abc", "def")
	e1.Call(e2.Id(), "NoSuchMethod", "abc", "def") // method not found
	e1.Call(e2.Id(), "Foo")                        // too few arguments
	e1.Call(e2.Id(), "Foo", 1, 2)                  // too many arguments
	time.Sleep(3 * time.Second)
}
