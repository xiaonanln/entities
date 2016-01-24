package sample_chatroom

import (
	"entities"
	"log"
)

type Avatar struct {
	entities.Entity
}

func (self *Avatar) Test(a float64, b float64, c string) {
	log.Printf("Avatar.Test %v %v %v", a, b, c)
}
