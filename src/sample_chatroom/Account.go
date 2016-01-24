package sample_chatroom

import (
	"entities"
	"log"
)

type Account struct {
	entities.Entity
}

func (self *Account) Login(user, pass string) {
	avatar, err := entities.NewEntity("Avatar")
	if err != nil {
		log.Printf("Create avatar failed: %s", err)
		return
	}

	log.Printf("Logining %s %s, avatar created: %v", user, pass, avatar)
	self.GiveClientTo(avatar)
}
