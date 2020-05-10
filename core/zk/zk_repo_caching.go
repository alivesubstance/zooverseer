package zk

import (
	"github.com/alivesubstance/zooverseer/core"
	goCache "github.com/goburrow/cache"
	goZk "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	// Zk path -> Node
	cache      goCache.Cache
	repository = Repository{}
)

type CachingRepositoryAccessor interface {
	InvalidateAll()
	Invalidate(zkPath string)

	Accessor
}

type CachingRepository struct {
	CachingRepositoryAccessor
}

func init() {
	stats := &goCache.Stats{}
	cache = goCache.New(goCache.WithExpireAfterAccess(core.ZkCacheExpireAfterAccessMinutes * time.Minute))
	// TODO looks like stats doesn't collect numbers
	cache.Stats(stats)

	go func() {
		time.Sleep(core.ZkCacheStatsPeriodMinutes * time.Minute)
		log.Infof("Zk cache: %+v\n", stats)
	}()
}

func (c *CachingRepository) SetConnInfo(connInfo *core.ConnInfo) {
	repository.SetConnInfo(connInfo)
}

func (c *CachingRepository) GetRootNode() (*Node, error) {
	var err error

	rootNode, _ := cache.GetIfPresent(core.NodeRootName)
	if rootNode == nil {
		rootNode, err = repository.GetRootNode()
		if err != nil {
			return nil, err
		}
		cache.Put(core.NodeRootName, rootNode)
	}

	return rootNode.(*Node), nil
}

func (c *CachingRepository) Get(path string) (*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		node, err = repository.Get(path)
		if node != nil {
			cache.Put(path, node)
		}
	}

	if err != nil {
		return nil, err
	}

	return node.(*Node), err
}

func (c *CachingRepository) GetMeta(path string) (*goZk.Stat, error) {
	var err error
	var meta *goZk.Stat
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		meta, err = repository.GetMeta(path)
		if meta != nil {
			cache.Put(path, &Node{Meta: meta})
		}
	} else if node.(*Node).Meta == nil {
		meta, err = repository.GetMeta(path)
		if meta != nil {
			node.(*Node).Meta = meta
			cache.Put(path, node)
		}
	}

	return meta, err
}

func (c *CachingRepository) GetValue(path string) (*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		node, err = repository.GetValue(path)
		if node != nil {
			cache.Put(path, node)
		}
	} else if len(node.(*Node).Value) == 0 {
		valueNode, valueNodeErr := repository.GetValue(path)
		err = valueNodeErr
		if valueNode != nil {
			node.(*Node).Value = valueNode.Value
			node.(*Node).Meta = valueNode.Meta
			cache.Put(path, node)
		}
	}

	if err != nil {
		return nil, err
	}

	return node.(*Node), err
}

func (c *CachingRepository) GetChildren(path string) ([]*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		children, childrenErr := repository.GetChildren(path)
		err = childrenErr
		if children != nil {
			node = &Node{Children: children}
			cache.Put(path, node)
		}
	} else if node.(*Node).Children == nil {
		children, childrenErr := repository.GetChildren(path)
		err = childrenErr
		if children != nil {
			node.(*Node).Children = children
			cache.Put(path, node)
		}
	}

	if node == nil && err != nil {
		return nil, err
	}

	return node.(*Node).Children, err
}

func (c *CachingRepository) Invalidate(path string) {
	log.Tracef("Invalidate %s", path)
	cache.Invalidate(path)
}

func (c *CachingRepository) InvalidateAll() {
	cache.InvalidateAll()
}

func (c *CachingRepository) Save(parentPath string, childName string, acl []goZk.ACL) error {
	return repository.Save(parentPath, childName, acl)
}

func (c *CachingRepository) Delete(path string, version int32) error {
	c.Invalidate(path)
	return repository.Delete(path, version)
}
