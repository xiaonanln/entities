package entities

import (
	. "common"
)

type Client interface {
	Call(id Eid, method string, args ...interface{})
	Close() error
}
