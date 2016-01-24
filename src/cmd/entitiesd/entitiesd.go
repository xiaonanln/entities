package main

import (
	"entitiesd"
	"log"
	"sample_chatroom"
)

func main() {
	entitiesd.Prepare()

	entitiesd.RegisterEntity(&entitiesd.TestEntity{})
	entitiesd.RegisterEntity(&sample_chatroom.ChatRoomStub{}) // chatroom manager
	entitiesd.RegisterEntity(&sample_chatroom.ChatRoom{})     // chatroom
	entitiesd.RegisterEntity(&sample_chatroom.Account{})      // Account
	entitiesd.RegisterEntity(&sample_chatroom.Avatar{})       // Avatar

	chatroomStub, err := entitiesd.NewGlobalEntity("ChatRoomStub")
	for err != nil {
		log.Fatalf("Creating global [%s] failed: %s", "ChatRoomStub", err)
		chatroomStub, err = entitiesd.NewGlobalEntity("ChatRoomStub") // re-create global ChatRoomStub if fail
	}

	log.Printf("Global ChatRoomStub created: %s, starting entitiesd ...", chatroomStub)
	entitiesd.Run()
}
