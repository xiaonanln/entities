package main

import (
	"conf"
	"fmt"
)

func main() {
	fmt.Println("GetExeFullPath:", conf.GetExeFullPath())
	fmt.Println("GetEntitiesRoot:", conf.GetEntitiesRoot())
	fmt.Println("GetEntitiesConfigPath:", conf.GetEntitiesConfigPath())
	fmt.Printf("GetEntitiesConfig: \n%v\n", *conf.GetEntitiesConfig())
	// mapdConf := conf.GetMapdConfig()
	// fmt.Println("mapd config:\n", mapdConf)
}
