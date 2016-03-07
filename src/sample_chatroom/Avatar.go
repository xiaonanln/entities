package sample_chatroom

import (
	"entities"
	"log"
)

type Avatar struct {
	entities.Entity
}

func (self *Avatar) OnGetClient() {
	log.Printf("Avatar %s OnGetClient %s", self, self.Client)
	self.CallGlobalEntity("ChatRoomStub", "OnAvatarGetClient", self.Id())
}

func (self *Avatar) OnLoseClient(oldClient *entities.Client) {
	log.Printf("Avatar %s OnLoseClient", self)
	self.CallGlobalEntity("ChatRoomStub", "OnAvatarLoseClient", self.Id())
}

func (self *Avatar) Test(a float64, b float64, c string) {
	log.Printf("Avatar.Test %v %v %v", a, b, c)
	self.Call(self.Id(), "Test2")
}

func (self *Avatar) QueryChatRoomList() {
	log.Printf("Avatar.QueryChatRoomList ...")
	self.CallGlobalEntity("ChatRoomStub", "QueryChatRoomList", self.Id())
}

func (self *Avatar) OnQueryChatRoomList(roomIds []int) {
	log.Printf("Avatar.OnQueryChatRoomList: rooms = %v", roomIds)
}
