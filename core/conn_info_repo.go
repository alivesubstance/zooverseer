package core

import (
	"encoding/json"
	"fmt"
	"github.com/alivesubstance/zooverseer/util"
	"io/ioutil"
	"os"
)

//TODO remove it and put conn_info_repo.go into conn_repo folder?
type ConnRepository interface {
	Upsert(connInfo *JsonConnInfo)
	Find(connName string) (*JsonConnInfo, bool)
	FindAll() []JsonConnInfo
	Delete(connName string)
}

type ZooverseerConfig struct {
	Connections []JsonConnInfo /* `json:"connections"`*/
}

//TODO replace JsonConnInfo with JsonConnRepository
//type JsonConnRepository struct {
//	ConnRepository
//}

type JsonConnInfo struct {
	Name     string /*`json:"name"`*/
	Host     string
	Port     uint16
	User     string
	Password string
}

func (c *JsonConnInfo) String() string {
	return fmt.Sprintf(
		"JsonConnInfo[name: %s, host: %s, port: %v, user: %v]",
		c.Name, c.Host, c.Port, c.User,
	)
}

func (c *JsonConnInfo) Upsert(connInfo *JsonConnInfo) {

}

func (c *JsonConnInfo) Find(connName string) (*JsonConnInfo, bool) {
	if len(connName) == 0 {
		return nil, false
	}

	// can be replaced with json path but it also need fully read json file
	for _, connInfo := range c.FindAll() {
		if connInfo.Name == connName {
			return &connInfo, true
		}
	}

	return nil, false
}

func (c *JsonConnInfo) FindAll() []JsonConnInfo {
	config := readConfig()
	return config.Connections
}

func (c *JsonConnInfo) Delete(connName string) {

}

func readConfig() *ZooverseerConfig {
	var config ZooverseerConfig

	connConfigJson, err := os.Open(ConnConfigFilePath)
	util.CheckError(err)

	connConfigBytes, err := ioutil.ReadAll(connConfigJson)
	util.CheckError(err)

	json.Unmarshal(connConfigBytes, &config)
	defer connConfigJson.Close()

	return &config
}
