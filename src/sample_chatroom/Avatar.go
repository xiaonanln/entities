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
	self.Call(self.Id(), "Test2")
}
func (self *Avatar) Test2() {
	log.Println("Avatar.Test2 called from Test")
}

func (self *Avatar) OnGetClient() {

}

func (self *Avatar) OnLoseClient(oldClient *entities.Client) {

}
