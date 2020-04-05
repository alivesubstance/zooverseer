package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/avast/retry-go"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	gopath "path"
	"time"
)

type Node struct {
	Name     string
	Value    string
	Meta     *zkGo.Stat
	Children []Node
}

var retryOptions = []retry.Option{
	retry.Attempts(core.ZkOpRetryAttempts),
	retry.Delay(core.ZkOpRetryDelay * time.Millisecond),
	retry.OnRetry(func(n uint, err error) {
		log.WithError(err).Infof("Zk op failed. Retry %v of %v", n, core.ZkOpRetryAttempts)
	}),
}

func GetRootNodeChildren(connInfo *core.JsonConnInfo) ([]Node, error) {
	childPathCreator := func(path string, childName string) string {
		return fmt.Sprintf("/%s", childName)
	}

	return doGetChildren(core.NodeRootName, connInfo, childPathCreator)
}

func Get(path string, connInfo *core.JsonConnInfo) (*Node, error) {
	log.Info("Get data for " + path)

	value, meta, err := GetValue(path, connInfo)
	if err != nil {
		return nil, err
	}

	children, err := GetChildren(path, connInfo)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Name:     gopath.Base(path),
		Value:    value,
		Meta:     meta,
		Children: children,
	}

	return node, nil
}

func Exists(path string, connInfo *core.JsonConnInfo) (bool, *zkGo.Stat, error) {
	conn, err := getConn(connInfo)
	if err != nil {
		log.WithError(err).Errorf("Failed to check existing for %s", path)
		return false, nil, err
	}

	return conn.Exists(path)
}

func GetValue(path string, connInfo *core.JsonConnInfo) (string, *zkGo.Stat, error) {
	log.Info("Looking for value for " + path)

	var value string
	var meta *zkGo.Stat

	conn, err := getConn(connInfo)
	if err != nil {
		return "", nil, err
	}

	err = retry.Do(
		func() error {
			valueBytes, stat, err := conn.Get(path)
			if err != nil {
				log.WithError(err).Error("Failed to get value for " + path)
				return err
			}

			value = util.BytesToString(valueBytes)
			meta = stat
			return nil
		}, retryOptions...,
	)

	if err != nil {
		return "", nil, err
	}

	return value, meta, nil
}

func GetChildren(path string, connInfo *core.JsonConnInfo) ([]Node, error) {
	childPathCreator := func(path string, childName string) string {
		return fmt.Sprintf("%s/%s", path, childName)
	}
	return doGetChildren(path, connInfo, childPathCreator)
}

func doGetChildren(path string, connInfo *core.JsonConnInfo, childPathCreator func(path string, childName string) string) ([]Node, error) {
	log.Info("Looking for children for " + path)

	conn, err := getConn(connInfo)
	if err != nil {
		return nil, err
	}

	//TODO add retry
	childrenNames, _, err := conn.Children(path)
	if err != nil {
		log.WithError(err).Fatal("Failed to get children for " + path)
		return nil, err
	}

	if len(childrenNames) == 0 {
		return nil, nil
	}

	nodes := make([]Node, len(childrenNames))
	for i, childName := range childrenNames {
		_, meta, err := Exists(childPathCreator(path, childName), connInfo)
		if err != nil {
			return nil, err
		}
		nodes[i] = Node{
			Name: childName,
			Meta: meta,
		}
	}

	return nodes, nil
}

func getConn(connInfo *core.JsonConnInfo) (*zkGo.Conn, error) {
	// dereferencing conn info to use struct copy(not a pointer) as a cache key
	connInfoValue := *connInfo
	conn, err := ConnCache.Get(connInfoValue)
	if err != nil {
		log.WithError(err).Errorf("Failed to get connection for %s", connInfo.String())
		return nil, err
	}

	return conn.(*zkGo.Conn), nil
}
