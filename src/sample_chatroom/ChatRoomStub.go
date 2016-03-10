package sample_chatroom

import (
	. "common"
	"entities"
	"log"
)

type ChatRoomStub struct {
	entities.Entity
	chatrooms      map[int]Eid
	nextChatRoomId int
}

func (self *ChatRoomStub) OnInit() {
	log.Printf("######################## ChatRoomStub.OnInit ########################")
	self.chatrooms = make(map[int]Eid)
	self.nextChatRoomId = 1000

	self.newChatRoom()
}

func (self *ChatRoomStub) OnAvatarGetClient(eid Eid) {
	log.Printf("ChatRoomStub.OnAvatarGetClient %v", eid)
}

func (self *ChatRoomStub) OnAvatarLoseClient(eid Eid) {
	log.Printf("ChatRoomStub.OnAvatarLoseClient %v", eid)
}

func (self *ChatRoomStub) OnRegisteredGlobally() {
	log.Printf("Register global entity: %s", self)
}

func (self *ChatRoomStub) QueryChatRoomList(caller string) {
	roomIds := []int{}
	for roomId := range self.chatrooms {
		roomIds = append(roomIds, roomId)
	}
	log.Printf("QueryChatRoomList: caller = %s", caller)
	callerEid := Eid(caller)
	self.Call(callerEid, "OnQueryChatRoomList", roomIds)
}

func (self *ChatRoomStub) newChatRoom() Eid {
	entity, err := self.NewEntity("ChatRoom")
	if err != nil {
		return ""
	}

	roomId := self.nextChatRoomId
	self.nextChatRoomId = roomId + 1
	eid := entity.Id()
	self.chatrooms[roomId] = eid
	return eid
}
