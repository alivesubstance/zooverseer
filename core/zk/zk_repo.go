package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/avast/retry-go"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	gopath "path"
	"sort"
	"time"
)

type Node struct {
	Name     string
	Value    string
	Meta     *zkGo.Stat
	Children []*Node
}

var retryOptions = []retry.Option{
	retry.Attempts(core.ZkOpRetryAttempts),
	retry.Delay(core.ZkOpRetryDelay * time.Millisecond),
	retry.OnRetry(func(n uint, err error) {
		log.WithError(err).Infof("Zk op failed. Retry %v of %v", n, core.ZkOpRetryAttempts)
	}),
}

type Accessor interface {
	Get(path string, connInfo *core.JsonConnInfo) (*Node, error)
	GetMeta(path string, connInfo *core.JsonConnInfo) (*zkGo.Stat, error)
	// Returns node value and metadata
	GetValue(path string, connInfo *core.JsonConnInfo) (*Node, error)
	GetChildren(path string, connInfo *core.JsonConnInfo) ([]*Node, error)
	GetRootNodeChildren(connInfo *core.JsonConnInfo) ([]*Node, error)
}

type Repository struct {
	Accessor
}

func (repo *Repository) GetRootNodeChildren(connInfo *core.JsonConnInfo) ([]*Node, error) {
	childPathCreator := func(path string, childName string) string {
		return fmt.Sprintf("/%s", childName)
	}

	return doGetChildren(nil, core.NodeRootName, connInfo, childPathCreator)
}

func (repo *Repository) Get(path string, connInfo *core.JsonConnInfo) (*Node, error) {
	log.Info("Get data for " + path)

	node, err := repo.GetValue(path, connInfo)
	if err != nil {
		return nil, err
	}

	children, err := repo.GetChildren(path, connInfo)
	if err != nil {
		return nil, err
	}

	node.Name = gopath.Base(path)
	node.Children = children

	return node, nil
}

func (repo *Repository) GetMeta(path string, connInfo *core.JsonConnInfo) (*zkGo.Stat, error) {
	conn, err := getConn(connInfo)
	if err != nil {
		log.WithError(err).Errorf("Failed to check existing for %s", path)
		return nil, err
	}

	var meta *zkGo.Stat
	err = retry.Do(
		func() error {
			_, meta, err = conn.Exists(path)
			if err != nil {
				log.WithError(err).Error("Failed to get metadata for " + path)
				return err
			}

			return nil
		}, retryOptions...)

	return meta, err
}

func (repo *Repository) GetValue(path string, connInfo *core.JsonConnInfo) (*Node, error) {
	log.Debug("Looking for value for " + path)

	conn, err := getConn(connInfo)
	if err != nil {
		return nil, err
	}

	var value []byte
	var meta *zkGo.Stat
	err = retry.Do(
		func() error {
			value, meta, err = conn.Get(path)
			if err != nil {
				log.WithError(err).Error("Failed to get value for " + path)
				return err
			}

			return nil
		}, retryOptions...,
	)

	if err != nil {
		return nil, err
	}

	node := &Node{
		Value: util.BytesToString(value),
		Meta:  meta,
	}
	return node, nil
}

func (repo *Repository) GetChildren(path string, connInfo *core.JsonConnInfo) ([]*Node, error) {
	childPathCreator := func(path string, childName string) string {
		return fmt.Sprintf("%s/%s", path, childName)
	}
	return doGetChildren(repo, path, connInfo, childPathCreator)
}

func doGetChildren(
	zkRepo *Repository, path string, connInfo *core.JsonConnInfo, childPathCreator func(path string, childName string) string,
) ([]*Node, error) {
	log.Info("Looking for children for " + path)

	conn, err := getConn(connInfo)
	if err != nil {
		return nil, err
	}

	var childrenNames []string
	err = retry.Do(
		func() error {
			childrenNames, _, err = conn.Children(path)
			if err != nil {
				log.WithError(err).Fatal("Failed to get children for " + path)
				return err
			}
			return nil
		}, retryOptions...,
	)

	if len(childrenNames) == 0 {
		return nil, nil
	}

	// get metadata for each child. mostly to know is there are children in child node
	nodes := make([]*Node, len(childrenNames))
	for i, childName := range childrenNames {
		meta, err := zkRepo.GetMeta(childPathCreator(path, childName), connInfo)
		if err != nil {
			return nil, err
		}
		nodes[i] = &Node{
			Name: childName,
			Meta: meta,
		}
	}

	//TODO make nodes with children first in the slice
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

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
