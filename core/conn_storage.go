package core

import (
	"encoding/json"
	"github.com/alivesubstance/zooverseer/util"
	"io/ioutil"
	"os"
)

type ConnRepository interface {
	Upset(connInfo JsonConnInfo)
	Find(connName string) (*JsonConnInfo, bool)
	FindAll() []JsonConnInfo
	Delete(connName string)
}

type ZooverseerConfig struct {
	Connections []JsonConnInfo /* `json:"connections"`*/
}

type JsonConnInfo struct {
	Name     string /*`json:"name"`*/
	Host     string
	Port     uint16
	User     string
	Password string
}

func (c JsonConnInfo) Upset(connInfo JsonConnInfo) {

}

func (c JsonConnInfo) Find(connName string) (*JsonConnInfo, bool) {
	if len(connName) == 0 {
		return nil, false
	}

	for _, connInfo := range c.FindAll() {
		if connInfo.Name == connName {
			return &connInfo, true
		}
	}

	return nil, false
}

func (c JsonConnInfo) FindAll() []JsonConnInfo {
	config := readConfig()
	return config.Connections
}

func (c JsonConnInfo) Delete(connName string) {

}

func readConfig() ZooverseerConfig {
	var config ZooverseerConfig

	connConfigJson, err := os.Open(ConnConfigFilePath)
	util.CheckError(err)

	connConfigBytes, err := ioutil.ReadAll(connConfigJson)
	util.CheckError(err)

	json.Unmarshal(connConfigBytes, &config)
	defer connConfigJson.Close()

	return config
}
