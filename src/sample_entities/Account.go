package sample_entities

import (
	"entities"
	"log"
)

type Account struct {
	entities.Entity
}

func (self *Account) Login(user, pass string) {
	avatar := entities.NewEntity("Avatar")
	log.Printf("Logining %s %s, avatar created: %v", user, pass, avatar)
	self.GiveClientTo(avatar)
}
