package core

import (
	"encoding/json"
	"fmt"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

//TODO remove it and put conn_info_repo.go into conn_repo folder?
type ConnRepository interface {
	Upsert(connInfo *ConnInfo)
	Find(connName string) *ConnInfo
	FindAll() []*ConnInfo
	Delete(connInfo *ConnInfo)
}

type JsonConnRepository struct {
	ConnRepository
}

type ConnInfo struct {
	Name     string /*`json:"name"`*/
	Host     string
	Port     int
	User     string
	Password string
}

func (c *ConnInfo) String() string {
	return fmt.Sprintf("%s(%s:%d)", c.Name, c.Host, c.Port)
}

func (c *ConnInfo) Copy() *ConnInfo {
	connCopy := &ConnInfo{}
	err := copier.Copy(&connCopy, c)
	if err != nil {
		log.WithError(err).Errorf("Copy")
	}

	return connCopy
}

func (c *JsonConnRepository) Upsert(connInfo *ConnInfo) {
	foundConnInfo := c.Find(connInfo.Name)
	if foundConnInfo != nil {
		c.Delete(foundConnInfo)
		c.insert(connInfo)
	} else {
		c.insert(connInfo)
	}
}

func (c *JsonConnRepository) insert(connInfo *ConnInfo) {
	conns := append(c.FindAll(), connInfo)
	c.saveAll(conns)
}

func (c *JsonConnRepository) Find(connName string) *ConnInfo {
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

func (c *JsonConnRepository) FindAll() []*ConnInfo {
	return readConns()
}

func (c *JsonConnRepository) Delete(connInfo *ConnInfo) {
	conns := readConns()
	if len(conns) == 0 {
		log.Warnf("There are no saved connections. Nothing to remove")
		return
	}

	var newConns = make([]*ConnInfo, 0)
	for _, conn := range conns {
		if conn.Name == connInfo.Name {
			continue
		}
		newConns = append(newConns, conn)
	}

	c.saveAll(newConns)
}

func (c *JsonConnRepository) saveAll(connInfos []*ConnInfo) {
	connConfigJson, err := os.Open(ConnConfigFilePath)
	util.CheckError(err)
	defer connConfigJson.Close()

	bytes, err := json.Marshal(connInfos)
	if err != nil {
		log.WithError(err).Errorf("Failed to marshal connections")
	}
	err = ioutil.WriteFile(ConnConfigFilePath, bytes, 0644)
	if err != nil {
		log.WithError(err).Errorf("Failed to write connections config")
	}
}

func readConns() []*ConnInfo {
	connInfos := make([]*ConnInfo, 0)

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
