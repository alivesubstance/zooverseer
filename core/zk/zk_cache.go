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
	stats := goCache.Stats{}
	cache = goCache.New(goCache.WithExpireAfterAccess(core.ZkCacheExpireAfterAccessMinutes * time.Minute))
	cache.Stats(&stats)

	go func() {
		time.Sleep(core.ZkCacheStatsPeriodMinutes * time.Minute)
		log.Infof("Zk cache: %+v\n", stats)
	}()
}

func (c *CachingRepository) GetRootNodeChildren(connInfo *core.JsonConnInfo) ([]*Node, error) {
	var err error
	var node Node

	children, _ := cache.GetIfPresent(core.NodeRootName)
	if children == nil {
		children, err = repository.GetRootNodeChildren(connInfo)
		if children != nil {
			node = Node{Name: core.NodeRootName, Children: children.([]*Node)}
			cache.Put(core.NodeRootName, &node)
		}
	}

	if err != nil {
		return nil, err
	}

	return children.([]*Node), nil
}

func (c *CachingRepository) Get(path string, connInfo *core.JsonConnInfo) (*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		node, err = repository.Get(path, connInfo)
		if node != nil {
			cache.Put(path, node)
		}
	}

	if err != nil {
		return nil, err
	}

	return node.(*Node), err
}

func (c *CachingRepository) GetMeta(path string, connInfo *core.JsonConnInfo) (*goZk.Stat, error) {
	var err error
	var meta *goZk.Stat
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		meta, err = repository.GetMeta(path, connInfo)
		if meta != nil {
			cache.Put(path, &Node{Meta: meta})
		}
	} else if node.(*Node).Meta == nil {
		meta, err = repository.GetMeta(path, connInfo)
		if meta != nil {
			node.(*Node).Meta = meta
			cache.Put(path, node)
		}
	}

	return meta, err
}

func (c *CachingRepository) GetValue(path string, connInfo *core.JsonConnInfo) (*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		node, err = repository.GetValue(path, connInfo)
		if node != nil {
			cache.Put(path, node)
		}
	} else if len(node.(*Node).Value) == 0 {
		valueNode, valueNodeErr := repository.GetValue(path, connInfo)
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

func (c *CachingRepository) GetChildren(path string, connInfo *core.JsonConnInfo) ([]*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		children, childrenErr := repository.GetChildren(path, connInfo)
		err = childrenErr
		if children != nil {
			node = &Node{Children: children}
			cache.Put(path, node)
		}
	} else if node.(*Node).Children == nil {
		children, childrenErr := repository.GetChildren(path, connInfo)
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
	cache.Invalidate(path)
}

func (c *CachingRepository) InvalidateAll() {
	cache.InvalidateAll()
}
