package entitiesd

import (
	"entities"
	"fmt"
)

type TestEntity struct {
	entities.Entity
}

func (self *TestEntity) Foo(a int) {
	fmt.Println("RPC SUCCESSFULLY:", self, "Foo", a)
}

func (self *TestEntity) Bar(a string, b string) {
	fmt.Println("RPC:", self, "Bar", a, b)
}
