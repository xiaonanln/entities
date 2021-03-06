package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	entitiesConfig *EntitiesConfig
)

type MapdConfig struct {
	Host    string
	CmdPort int `json:"cmd_port"`
	RPCPort int `json:"rpc_port"`
}

type GatedConfig struct {
	Host string
	Port int
}

type EntitiesdConfig struct {
	Host string
	Port int
}

type EntitiesConfig struct {
	BootEntity string `json:"boot_entity"`
	Mapd       MapdConfig
	Entitiesd  []EntitiesdConfig
	Gated      []GatedConfig
}

func GetExeFullPath() string {
	exePath := os.Args[0]
	exePath, err := filepath.Abs(exePath)
	if err != nil {
		panic(err)
	}
	return exePath
}

func GetEntitiesRoot() string {
	exePath := GetExeFullPath()
	exeDirPath := filepath.Dir(exePath)
	bin, entitiesRoot := filepath.Base(exeDirPath), filepath.Dir(exeDirPath)
	if bin != "bin" {
		panic(fmt.Errorf("executive not in entities/bin"))
	}
	return entitiesRoot
}

func GetEntitiesConfigPath() string {
	entitiesRoot := GetEntitiesRoot()
	return filepath.Join(entitiesRoot, "src", "conf", "entities.conf")
}

func GetEntitiesConfig() *EntitiesConfig {
	if entitiesConfig != nil {
		return entitiesConfig
	}

	configPath := GetEntitiesConfigPath()
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config EntitiesConfig
	// fmt.Println("read config:\n", string(content))
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	entitiesConfig = &config // cache the result
	return entitiesConfig
}

func GetMapdConfig() *MapdConfig {
	entitiesConfig := GetEntitiesConfig()
	return &entitiesConfig.Mapd
}

func GetBootEntity() string {
	return GetEntitiesConfig().BootEntity
}

func GetEntitiesdConfig(pid int) *EntitiesdConfig {
	config := GetEntitiesConfig()
	return &config.Entitiesd[pid-1]
}

func GetGatedConfig(gid int) *GatedConfig {
	config := GetEntitiesConfig()
	return &config.Gated[gid-1]
}

func GetEntitiesdCount() int {
	config := GetEntitiesConfig()
	return len(config.Entitiesd)
}

func GetGatedCount() int {
	config := GetEntitiesConfig()
	return len(config.Gated)
}
