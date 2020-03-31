package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	zookeeper "github.com/outbrain/zookeepercli/go/zk"
	log "github.com/sirupsen/logrus"
	gopath "path"
)

var zk = zookeeper.NewZooKeeper()

type Node struct {
	Name     string
	Value    string
	Children []Node
}

func Get(path string, connInfo *core.JsonConnInfo, chanNode chan Node) error {
	zk.SetServers(getServer(connInfo))
	if len(connInfo.User) != 0 && len(connInfo.Password) != 0 {
		authExp := fmt.Sprint(connInfo.User, ":", connInfo.Password)
		zk.SetAuth("digest", []byte(authExp))
	}

	children, err := getChildren(path)
	if err != nil {
		return err
	}

	//value, err := getValue(path)
	//if err != nil {
	//	return err
	//}

	node := Node{
		Name:     gopath.Base(path),
		Value:    "",
		Children: children,
	}

	chanNode <- node
	return nil
}

func getChildren(path string) ([]Node, error) {
	children, err := zk.Children(path)
	if err != nil {
		log.Error("Failed to get children for [" + path + "]")
		return nil, err
	}

	if len(children) == 0 {
		return nil, nil
	}

	nodes := make([]Node, len(children))
	for i, child := range children {
		nodes[i] = Node{Name: child}
	}

	return nodes, nil
}

func getValue(path string) (string, error) {
	log.Debug("Looking for value for [" + path + "]")
	valueBytes, err := zk.Get(path)
	if err != nil {
		log.Error("Failed to get value for [" + path + "]")
		return "", err
	}

	return util.BytesToString(valueBytes), nil
}

func getServer(info *core.JsonConnInfo) []string {
	servers := make([]string, 1)
	servers[0] = fmt.Sprintf("%v:%v", info.Host, info.Port)
	return servers
}
