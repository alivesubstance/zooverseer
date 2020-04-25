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
	Delete(connInfo *JsonConnInfo)
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
	return fmt.Sprintf("%s(%s:%d)", c.Name, c.Host, c.Port)
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
	return readConns()
}

func (c *JsonConnRepository) Delete(connInfo *JsonConnInfo) {
	conns := readConns()
	if len(conns) == 0 {
		log.Warnf("There are no saved connections. Nothing to remove")
		return
	}

	var newConns = make([]*JsonConnInfo, 0)
	for _, conn := range conns {
		if conn.Name == connInfo.Name {
			continue
		}
		newConns = append(newConns, conn)
	}

	connConfigJson, err := os.Open(ConnConfigFilePath)
	util.CheckError(err)
	defer connConfigJson.Close()

	bytes, err := json.Marshal(newConns)
	if err != nil {
		log.WithError(err).Errorf("Failed to marshal connections")
	}
	err = ioutil.WriteFile(ConnConfigFilePath, bytes, 0644)
	if err != nil {
		log.WithError(err).Errorf("Failed to write connections config")
	}
}

func readConns() []*JsonConnInfo {
	connInfos := make([]*JsonConnInfo, 0)

	connConfigJson, err := os.Open(ConnConfigFilePath)
	if os.IsNotExist(err) {
		log.Tracef("Connections file doesn't exist")
		return nil
	}
	util.CheckError(err)
	defer connConfigJson.Close()

	connConfigBytes, err := ioutil.ReadAll(connConfigJson)
	util.CheckError(err)

	err = json.Unmarshal(connConfigBytes, &connInfos)
	if err != nil {
		log.WithError(err).Errorf("Failed to read connections config")
	}

	return connInfos
}
