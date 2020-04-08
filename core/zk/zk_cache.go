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

type CachingRepository struct {
	Accessor
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

func (cachingRepo *CachingRepository) GetRootNodeChildren(connInfo *core.JsonConnInfo) ([]*Node, error) {
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

func (cachingRepo *CachingRepository) Get(path string, connInfo *core.JsonConnInfo) (*Node, error) {
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

func (cachingRepo *CachingRepository) GetMeta(path string, connInfo *core.JsonConnInfo) (*goZk.Stat, error) {
	var err error
	var meta *goZk.Stat
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		meta, err = repository.GetMeta(path, connInfo)
		if meta != nil {
			cache.Put(path, &Node{Meta: meta})
		}
	}

	return meta, err
}

func (cachingRepo *CachingRepository) GetValue(path string, connInfo *core.JsonConnInfo) (*Node, error) {
	var err error
	node, _ := cache.GetIfPresent(path)
	if node == nil {
		node, err = repository.GetValue(path, connInfo)
		if node != nil {
			cache.Put(path, node)
		}
	}

	if err != nil {
		return nil, err
	}

	return node.(*Node), err
}

func (cachingRepo *CachingRepository) GetChildren(path string, connInfo *core.JsonConnInfo) ([]Node, error) {
	var err error
	children, _ := cache.GetIfPresent(core.NodeRootName)
	if children == nil {
		children, err = repository.GetChildren(path, connInfo)
		if children != nil {
			cache.Put(path, &Node{Children: children.([]*Node)})
		}
	}

	if children == nil && err != nil {
		return nil, err
	}

	return children.([]Node), err
}

func Invalidate(path string) {
	cache.Invalidate(path)
}

func InvalidateAll() {
	cache.InvalidateAll()
}
