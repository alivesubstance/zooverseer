package core

import (
	"encoding/json"
	"github.com/alivesubstance/zooverseer/util"
	"io/ioutil"
	"os"
)

type ConnRepository interface {
	Upset(connInfo JsonConnInfo)
	Find(connName string) JsonConnInfo
	FindAll() []JsonConnInfo
	Delete(connName string)
}

func (c JsonConnInfo) Upset(connInfo JsonConnInfo) {

}

func (c JsonConnInfo) Find(connName string) JsonConnInfo {
	panic(nil)
}

func (c JsonConnInfo) FindAll() []JsonConnInfo {
	var connInfos JsonConnInfos

	connConfigJson, err := os.Open(ConnConfigFilePath)
	util.CheckError(err)

	connConfigBytes, err := ioutil.ReadAll(connConfigJson)
	util.CheckError(err)

	json.Unmarshal(connConfigBytes, &connInfos)
	defer connConfigJson.Close()

	return connInfos.Connections
}

func (c JsonConnInfo) Delete(connName string) {

}

type JsonConnInfos struct {
	Connections []JsonConnInfo /* `json:"connections"`*/
}

type JsonConnInfo struct {
	Name     string /*`json:"name"`*/
	Host     string
	Port     int16
	User     string
	Password string
}
