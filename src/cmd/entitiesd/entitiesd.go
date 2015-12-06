package main

import (
	"fmt"
	_ "rpc"
)

type Base struct {
	a int
	B int
}

func (self *Base) foo() {
	fmt.Println(self.a, self.B)
}

type Child struct {
	Base Base
}

func (self *Child) foo() {
	fmt.Println("child foo")
}

func main() {
	var b Base
	var c Child
	b.foo()

	c.foo()
	c.Base.foo()
	println(c.Base.B, c.Base.a)
}
