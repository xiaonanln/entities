package sample_chatroom

import (
	. "common"
	"entities"
	"log"
)

type ChatRoomStub struct {
	entities.Entity
}

func (self *ChatRoomStub) OnAvatarGetClient(eid Eid) {
	log.Printf("ChatRoomStub.OnAvatarGetClient %v", eid)
}

func (self *ChatRoomStub) OnAvatarLoseClient(eid Eid) {
	log.Printf("ChatRoomStub.OnAvatarLoseClient %v", eid)
}
