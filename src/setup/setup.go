package setup

import (
	"log"
	"os"
	"strings"
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

	// fixme : use better method
	if strings.HasSuffix(exePath, "entitiesd") || strings.HasSuffix(exePath, "entitiesd.exe") {
		role = ENTITIESD
	} else if strings.HasSuffix(exePath, "gated") || strings.HasSuffix(exePath, "gated.exe") {
		role = GATED
	} else if strings.HasSuffix(exePath, "mapd") || strings.HasSuffix(exePath, "mapd.exe") {
		role = MAPD
	} else {
	}

	log.Printf("Setup %s to be role %d.", exePath, role)
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
