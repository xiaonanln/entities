package entitiesd

import (
	. "common"
)

type Client struct {
	Gid      int
	ClientId ClientId
}

func NewClient(gid int, clientid ClientId) *Client {
	return &Client{
		Gid:      gid,
		ClientId: clientid,
	}
}

func (self *Client) Call(id Eid, method string, args ...interface{}) {

}

func (self *Client) Close() error {
	return nil
}
