package entitiesd

import (
	"conf"
	"log"
)

var (
	gateds []*EntitiesdClientProxy
)

func init() {
	gateds = make([]*EntitiesdClientProxy, conf.GetGatedCount())
	log.Printf("Found %d gated", len(gateds))
}

func onNewGated(gated *EntitiesdClientProxy, gid int) {
	if gateds[gid-1] != nil {
		gateds[gid-1].Close()
		gateds[gid-1] = nil
	}

	gated.SetGid(gid)
	gateds[gid-1] = gated
	log.Printf("New gated %s for gid %d", gated, gid)
}
