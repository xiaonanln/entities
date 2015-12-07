package main

import (
	"conf"
	"fmt"
)

func main() {
	fmt.Println("GetExeFullPath:", conf.GetExeFullPath())
	fmt.Println("GetEntitiesRoot:", conf.GetEntitiesRoot())
	fmt.Println("GetEntitiesConfigPath:", conf.GetEntitiesConfigPath())
	fmt.Printf("ReadEntitiesConfig: \n%v\n", *conf.ReadEntitiesConfig())
	// mapdConf := conf.ReadMapdConfig()
	// fmt.Println("mapd config:\n", mapdConf)
}
