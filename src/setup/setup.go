package setup

import (
	"log"
	"os"
	"path"
)

const (
	INVALID_ROLE = iota
	ENTITIESD    = iota
	GATED        = iota
	MAPD         = iota
)

var (
	role = INVALID_ROLE
)

func init() {
	exePath := os.Args[0]
	exeName := path.Base(exePath)
	if exeName == "entitiesd" {
		role = ENTITIESD
	} else if exeName == "gated" {
		role = GATED
	} else if exeName == "mapd" {
		role = MAPD
	} else {
	}

	log.Printf("Setup %s to be role %d.", exeName, role)
}

func GetRole() int {
	return role
}

func IsEntitiesd() bool {
	return role == ENTITIESD
}

func IsGated() bool {
	return role == GATED
}

func IsMapd() bool {
	return role == MAPD
}
