package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	gopath "path"
)

type Node struct {
	Name        string
	Value       string
	Meta        *zkGo.Stat
	HasChildren bool
	Children    []Node
}

func Get(path string, connInfo *core.JsonConnInfo /*, chanNode chan Node*/) (*Node, error) {
	log.Info("Get data for " + path)

	value, err := GetValue(path, connInfo)
	if err != nil {
		return nil, err
	}

	children, err := GetChildren(path, connInfo)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Name:        gopath.Base(path),
		Value:       value,
		HasChildren: len(children) > 0,
		Children:    children,
	}

	//chanNode <- node
	return node, nil
}

func Exists(path string, connInfo *core.JsonConnInfo) (bool, *zkGo.Stat, error) {
	conn, err := getConn(connInfo)
	if err != nil {
		log.Errorf("Failed to check existing for %s", path, err)
		return false, nil, err
	}

	return conn.Exists(path)
}

func GetValue(path string, connInfo *core.JsonConnInfo) (string, error) {
	log.Info("Looking for value for " + path)

	conn, err := getConn(connInfo)
	if err != nil {
		return "", err
	}

	valueBytes, _, err := conn.Get(path)
	if err != nil {
		log.Error("Failed to get value for " + path)
		return "", err
	}

	return util.BytesToString(valueBytes), nil
}

func GetChildren(path string, connInfo *core.JsonConnInfo) ([]Node, error) {
	log.Info("Looking for children for " + path)

	conn, err := getConn(connInfo)
	if err != nil {
		return nil, err
	}

	childrenNames, _, err := conn.Children(path)
	if err != nil {
		log.Error("Failed to get children for " + path)
		return nil, err
	}

	if len(childrenNames) == 0 {
		return nil, nil
	}

	nodes := make([]Node, len(childrenNames))
	for i, childName := range childrenNames {
		_, stat, err := Exists(fmt.Sprintf("%s/%s", path, childName), connInfo)
		if err != nil {
			return nil, err
		}
		nodes[i] = Node{
			Name:        childName,
			HasChildren: stat.NumChildren > 0,
		}
	}

	return nodes, nil
}

func getConn(connInfo *core.JsonConnInfo) (*zkGo.Conn, error) {
	conn, err := ConnCache.Get(connInfo.String())
	if err != nil {
		log.Errorf("Failed to get connection for %s", connInfo.String())
		return nil, err
	}

	return conn.(*zkGo.Conn), nil
}
