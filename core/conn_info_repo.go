package core

import (
	"encoding/json"
	"fmt"
	"github.com/alivesubstance/zooverseer/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

//TODO remove it and put conn_info_repo.go into conn_repo folder?
type ConnRepository interface {
	Upsert(connInfo *JsonConnInfo)
	Find(connName string) *JsonConnInfo
	FindAll() []*JsonConnInfo
	Delete(connName string)
}

type ZooverseerConfig struct {
	Connections []*JsonConnInfo /* `json:"connections"`*/
}

type JsonConnRepository struct {
	ConnRepository
}

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

func (c *JsonConnRepository) Upsert(connInfo *JsonConnInfo) {

}

func (c *JsonConnRepository) Find(connName string) *JsonConnInfo {
	if len(connName) == 0 {
		return nil
	}

	// can be replaced with json path but it also need fully read json file
	for _, connInfo := range c.FindAll() {
		if connInfo.Name == connName {
			return connInfo
		}
	}

	return nil
}

func (c *JsonConnRepository) FindAll() []*JsonConnInfo {
	config := readConfig()
	return config.Connections
}

func (c *JsonConnRepository) Delete(connName string) {

}

func readConfig() *ZooverseerConfig {
	var config ZooverseerConfig

	connConfigJson, err := os.Open(ConnConfigFilePath)
	util.CheckError(err)

	connConfigBytes, err := ioutil.ReadAll(connConfigJson)
	util.CheckError(err)

	err = json.Unmarshal(connConfigBytes, &config)
	if err != nil {
		log.WithError(err).Errorf("Failed to unmarshall Zooverseer config")
	}
	defer connConfigJson.Close()

	return &config
}
