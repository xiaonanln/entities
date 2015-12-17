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

type EntitiesdConfig struct {
	Host string
	Port int
	Pid  int
}

type EntitiesConfig struct {
	Mapd      MapdConfig
	Entitiesd []EntitiesdConfig
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

func ReadEntitiesConfig() *EntitiesConfig {
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

func ReadMapdConfig() *MapdConfig {
	entitiesConfig := ReadEntitiesConfig()
	return &entitiesConfig.Mapd
}
