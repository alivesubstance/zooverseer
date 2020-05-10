package zk

import (
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	goZk "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	gopath "path"
	"sync"
	"time"
)

var connCreateLock = sync.Mutex{}

type Node struct {
	Name     string
	Value    string
	Meta     *goZk.Stat
	Children []*Node
	Acl      []goZk.ACL
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
	SetConnInfo(connInfo *core.ConnInfo)
	Get(path string) (*Node, error)
	GetMeta(path string) (*goZk.Stat, error)
	// Returns node value and metadata
	GetValue(path string) (*Node, error)
	GetChildren(path string) ([]*Node, error)
	GetRootNode() (*Node, error)
	Save(parentPath string, childName string, acl []goZk.ACL) error
	Delete(path string, version int32) error
}

type Repository struct {
	Accessor

	connInfo *core.ConnInfo
}

func (r *Repository) Init(connInfo *core.ConnInfo) {
	r.connInfo = connInfo
}

func (r *Repository) SetConnInfo(connInfo *core.ConnInfo) {
	r.connInfo = connInfo
}

func (r *Repository) GetRootNode() (*Node, error) {
	var err error
	rootNode, err := r.GetValue(core.NodeRootName)
	children, err := r.GetChildren(core.NodeRootName)
	if err != nil {
		return nil, err
	}

	rootNode.Children = children
	return rootNode, nil
}

func (r *Repository) Get(path string) (*Node, error) {
	node, err := r.GetValue(path)
	if err != nil {
		return nil, err
	}

	children, err := r.GetChildren(path)
	if err != nil {
		return nil, err
	}

	node.Name = gopath.Base(path)
	node.Children = children

	return node, nil
}

func (r *Repository) GetMeta(path string) (*goZk.Stat, error) {
	conn, err := r.getConn()
	if err != nil {
		log.WithError(err).Errorf("Failed to check existing for %s", path)
		return nil, err
	}

	var meta *goZk.Stat
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

func (r *Repository) GetValue(path string) (*Node, error) {
	conn, err := r.getConn()
	if err != nil {
		return nil, err
	}

	var value []byte
	var meta *goZk.Stat
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
		Name:  gopath.Base(path),
		Value: util.BytesToString(value),
		Meta:  meta,
	}
	return node, nil
}

func (r *Repository) GetChildren(path string) ([]*Node, error) {
	conn, err := r.getConn()
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
	return r.getChildrenMeta(path, childrenNames)
}

func (r *Repository) Save(parentPath string, childName string, acl []goZk.ACL) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}

	absPath := r.buildAbsPath(parentPath, childName)
	_, err = conn.Create(absPath, util.StringToBytes(""), int32(0), acl)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(path string, version int32) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}

	return conn.Delete(path, version)
}

func (r *Repository) buildAbsPath(parentPath string, childName string) string {
	if parentPath == core.NodeRootName {
		return "/" + childName
	}
	return parentPath + "/" + childName
}

func (r *Repository) getChildrenMeta(path string, childrenNames []string) ([]*Node, error) {
	var resultErr error
	var wg sync.WaitGroup
	nodes := make([]*Node, len(childrenNames))
	for i, childName := range childrenNames {
		wg.Add(1)
		go func(idx int, path string, childName string) {
			defer wg.Done()

			absPath := r.buildAbsPath(path, childName)
			meta, err := r.GetMeta(absPath)
			if err != nil {
				log.WithError(err).Errorf("Failed to get node meta %s", absPath)
				resultErr = err
			}

			nodes[idx] = &Node{
				Name: childName,
				Meta: meta,
			}
		}(i, path, childName)
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

func (r *Repository) getConn() (*goZk.Conn, error) {
	//todo stupid cache doesn't lock loader function when call it after it didn't find entry it cache
	connCreateLock.Lock()
	defer connCreateLock.Unlock()

	// dereferencing conn info to use struct copy(not a pointer) as a cache key
	connInfoValue := *r.connInfo
	conn, err := ConnCache.Get(connInfoValue)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get connection for %s", r.connInfo.String())
	}

	return conn.(*goZk.Conn), nil
}
