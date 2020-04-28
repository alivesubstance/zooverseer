package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	gopath "path"
	"sync"
	"time"
)

var connCreateLock = sync.Mutex{}

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

//TODO measure timings for each operation
type Accessor interface {
	Get(path string, connInfo *core.ConnInfo) (*Node, error)
	GetMeta(path string, connInfo *core.ConnInfo) (*zkGo.Stat, error)
	// Returns node value and metadata
	GetValue(path string, connInfo *core.ConnInfo) (*Node, error)
	GetChildren(path string, connInfo *core.ConnInfo) ([]*Node, error)
	GetRootNodeChildren(connInfo *core.ConnInfo) ([]*Node, error)
}

type Repository struct {
	Accessor
}

func (repo *Repository) GetRootNodeChildren(connInfo *core.ConnInfo) ([]*Node, error) {
	absolutePathCreator := func(path string, childName string) string {
		return fmt.Sprintf("/%s", childName)
	}

	return doGetChildren(nil, core.NodeRootName, connInfo, absolutePathCreator)
}

func (repo *Repository) Get(path string, connInfo *core.ConnInfo) (*Node, error) {
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

func (repo *Repository) GetMeta(path string, connInfo *core.ConnInfo) (*zkGo.Stat, error) {
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

func (repo *Repository) GetValue(path string, connInfo *core.ConnInfo) (*Node, error) {
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

func (repo *Repository) GetChildren(path string, connInfo *core.ConnInfo) ([]*Node, error) {
	absolutePathCreator := func(path string, childName string) string {
		return fmt.Sprintf("%s/%s", path, childName)
	}
	return doGetChildren(repo, path, connInfo, absolutePathCreator)
}

func doGetChildren(
	zkRepo *Repository, path string, connInfo *core.ConnInfo, absolutePathCreator func(path string, childName string) string,
) ([]*Node, error) {
	log.Tracef("Looking for children for %s", path)

	conn, err := getConn(connInfo)
	if err != nil {
		return nil, err
	}

	var childrenNames []string
	err = retry.Do(
		func() error {
			childrenNames, _, err = conn.Children(path)
			if err != nil {
				log.WithError(err).Fatalf("Failed to get children for %s", path)
				return err
			}
			return nil
		}, retryOptions...,
	)

	if len(childrenNames) == 0 {
		return nil, nil
	}

	// get metadata for each child. mostly to know is there are children in child node
	return getChildrenMeta(zkRepo, path, connInfo, absolutePathCreator, childrenNames)
}

func getChildrenMeta(
	zkRepo *Repository,
	path string,
	connInfo *core.ConnInfo,
	absolutePathCreator func(path string, childName string) string,
	childrenNames []string,
) ([]*Node, error) {
	var resultErr error
	var wg sync.WaitGroup
	nodes := make([]*Node, len(childrenNames))
	for i, childName := range childrenNames {
		wg.Add(1)
		go func(idx int, path string, childName string, connInfo *core.ConnInfo) {
			defer wg.Done()

			absolutePath := absolutePathCreator(path, childName)
			meta, err := zkRepo.GetMeta(absolutePath, connInfo)
			if err != nil {
				log.WithError(err).Errorf("Failed to get node meta %s", absolutePath)
				resultErr = err
			}

			nodes[idx] = &Node{
				Name: childName,
				Meta: meta,
			}
		}(i, path, childName, connInfo)
	}
	wg.Wait()

	return sortNodes(nodes), resultErr
}

func sortNodes(nodes []*Node) []*Node {
	byName := func(n1, n2 *Node) bool { return n1.Name < n2.Name }
	byChildren := func(n1, n2 *Node) bool { return n1.Meta.NumChildren > 0 && n2.Meta.NumChildren <= 0 }
	var lessFuncs = []lessFunc{byName}

	if core.Config.SortFolderFirst {
		lessFuncs = []lessFunc{byChildren, byName}
	}

	OrderedBy(lessFuncs...).Sort(nodes)

	return nodes
}

func getConn(connInfo *core.ConnInfo) (*zkGo.Conn, error) {
	//todo stupid cache doesn't lock loader function when call it after it didn't find entry it cache
	connCreateLock.Lock()
	defer connCreateLock.Unlock()

	// dereferencing conn info to use struct copy(not a pointer) as a cache key
	connInfoValue := *connInfo
	conn, err := ConnCache.Get(connInfoValue)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get connection for %s", connInfo.String())
	}

	return conn.(*zkGo.Conn), nil
}
