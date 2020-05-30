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
	Name     string     `json:"name"`
	Value    string     `json:"value,omitempty"`
	Meta     *goZk.Stat `json:"-"`
	Children []*Node    `json:"children,omitempty"`
	Acl      []goZk.ACL `json:"-"`
}

var retryOptions = []retry.Option{
	retry.Attempts(core.ZkOpRetryAttempts),
	retry.Delay(core.ZkOpRetryDelay * time.Millisecond),
	retry.OnRetry(func(n uint, err error) {
		log.WithError(err).Infof("Zk op failed. Retry %v of %v", n, core.ZkOpRetryAttempts)
	}),
}

//todo measure timings for each operation
type Accessor interface {
	SetConnInfo(connInfo *core.ConnInfo)
	Get(path string) (*Node, error)

	//todo consider to remove next three methods and use Get instead
	GetMeta(path string) (*goZk.Stat, error)
	GetValue(path string) (*Node, error)
	GetChildren(path string) ([]*Node, error)

	GetRootNode() (*Node, error)
	SaveChild(parentPath string, childName string, acl []goZk.ACL) error
	SaveValue(parentPath string, node *Node) error
	Delete(path string, version int32) error
	Close()
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
				log.WithError(err).Errorf("Failed to get metadata for %s", path)
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
				log.WithError(err).Errorf("Failed to get value for %s", path)
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
				log.WithError(err).Error("Failed to get children for %s", path)
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

func (r *Repository) SaveChild(parentPath string, childName string, acl []goZk.ACL) error {
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

func (r *Repository) SaveValue(absPath string, node *Node) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}

	_, err = conn.Set(absPath, util.StringToBytes(node.Value), node.Meta.Version)
	return err
}

//todo create func like doWithConn(func action()) and
//to hide
// 	conn, err := r.getConn()
//	if err != nil {
//		return err
//	}

func (r *Repository) Delete(path string, node *Node) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}

	if node.Meta == nil {
		meta, err := r.GetMeta(path)
		if err != nil {
			return err
		}
		node.Meta = meta
	}

	if node.Meta.NumChildren > 0 {
		children, err := r.GetChildren(path)
		if err != nil {
			return err
		}

		for _, child := range children {
			childPath := path + "/" + child.Name
			log.Tracef("Deleting child %s", childPath)
			err := r.Delete(childPath, node)
			if err != nil {
				return err
			}
		}
	}

	log.Tracef("Deleting %s", path)
	return conn.Delete(path, node.Meta.Version)
}

func (r *Repository) Close() {
	conn, _ := r.getConn()
	if conn != nil {
		conn.Close()
	}
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
