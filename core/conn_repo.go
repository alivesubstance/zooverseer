package core

import (
	"encoding/json"
	"fmt"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type Type string

const (
	ConnManual    Type = "manual"
	ConnGenerated Type = "generated"
)

type ConnRepository interface {
	Upsert(connInfo *ConnInfo)
	FindById(id int64) *ConnInfo
	FindByName(name string) *ConnInfo
	FindAll() []*ConnInfo
	Delete(connInfo *ConnInfo)
	SaveAll(connInfos []*ConnInfo)
}

type JsonConnRepository struct {
	ConnRepository
}

type ConnInfo struct {
	Id       int64
	Name     string /*`json:"name"`*/
	Host     string
	Port     int
	User     string
	Password string
	Type     Type
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

	connCopy.Id = 0
	return connCopy
}

func (c *JsonConnRepository) Upsert(connInfo *ConnInfo) {
	if connInfo.Id == 0 {
		connInfo.Id = time.Now().Unix()
	}
	foundConnInfo := c.FindById(connInfo.Id)
	if foundConnInfo != nil {
		c.Delete(foundConnInfo)
		c.insert(connInfo)
	} else {
		c.insert(connInfo)
	}
}

func (c *JsonConnRepository) insert(connInfo *ConnInfo) {
	conns := append(c.FindAll(), connInfo)
	c.SaveAll(conns)
}

func (c *JsonConnRepository) FindById(id int64) *ConnInfo {
	if id == 0 {
		return nil
	}

	// can be replaced with json path but it also need fully read json file
	for _, connInfo := range c.FindAll() {
		if connInfo.Id == id {
			if len(connInfo.Password) != 0 {
				connInfo.Password = util.Decrypt(connInfo.Password)
			}
			return connInfo
		}
	}

	return nil
}

func (c *JsonConnRepository) FindByName(name string) *ConnInfo {
	if len(name) == 0 {
		return nil
	}

	// can be replaced with json path but it also need fully read json file
	for _, connInfo := range c.FindAll() {
		if connInfo.Name == name {
			if len(connInfo.Password) != 0 {
				connInfo.Password = util.Decrypt(connInfo.Password)
			}
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
		if conn.Id == connInfo.Id {
			continue
		}
		newConns = append(newConns, conn)
	}

	c.SaveAll(newConns)
}

func (c *JsonConnRepository) SaveAll(connInfos []*ConnInfo) {
	var file *os.File
	var err error

	_, err = os.Stat(Config.ConnConfigFilePath)
	if os.IsNotExist(err) {
		file, err = os.Create(Config.ConnConfigFilePath)
		log.Infof("File %s has been created", Config.ConnConfigFilePath)
	} else {
		file, err = os.Open(Config.ConnConfigFilePath)
	}
	util.CheckError(err)
	defer file.Close()

	bytes, err := json.Marshal(connInfos)
	if err != nil {
		log.WithError(err).Errorf("Failed to marshal connections")
	}
	err = ioutil.WriteFile(Config.ConnConfigFilePath, bytes, 0644)
	if err != nil {
		log.WithError(err).Errorf("Failed to write connections config")
	}
}

func readConns() []*ConnInfo {
	connInfos := make([]*ConnInfo, 0)

	connConfigJson, err := os.Open(Config.ConnConfigFilePath)
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
