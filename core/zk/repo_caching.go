package zk

import (
	"context"
	"github.com/alivesubstance/zooverseer/core"
	goCache "github.com/goburrow/cache"
	goZk "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	// Zk path -> Data
	cache goCache.Cache
)

type CachingRepositoryAccessor interface {
	InvalidateAll()
	Invalidate(path string)
	Export(path string) (*Node, error)

	Accessor
}

type CachingRepository struct {
	CachingRepositoryAccessor

	Repo Repository
}

func init() {
	cache = goCache.New(goCache.WithExpireAfterAccess(core.ZkCacheExpireAfterAccessMinutes * time.Minute))

	go func() {
		for {
			time.Sleep(core.ZkCacheStatsPeriodMinutes * time.Minute)
			stats := goCache.Stats{}
			cache.Stats(&stats)
			log.Infof("Zk cache: %+v", stats)
		}
	}()
}

func (c *CachingRepository) SetConnInfo(connInfo *core.ConnInfo) {
	c.Repo.SetConnInfo(connInfo)
}

func (c *CachingRepository) GetRootNode() (*Node, error) {
	var err error

	rootNode, _ := cache.GetIfPresent(core.NodeRootName)
	if rootNode == nil {
		rootNode, err = c.Repo.GetRootNode()
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
		node, err = c.Repo.Get(path)
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
		meta, err = c.Repo.GetMeta(path)
		if meta != nil {
			cache.Put(path, &Node{Meta: meta})
		}
	} else if node.(*Node).Meta == nil {
		meta, err = c.Repo.GetMeta(path)
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
		node, err = c.Repo.GetValue(path)
		if node != nil {
			cache.Put(path, node)
		}
	} else if len(node.(*Node).Value) == 0 {
		valueNode, valueNodeErr := c.Repo.GetValue(path)
		err = valueNodeErr
		if valueNode != nil {
			node.(*Node).Value = valueNode.Value
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
		children, childrenErr := c.Repo.GetChildren(path)
		err = childrenErr
		if children != nil {
			node = &Node{Children: children}
			cache.Put(path, node)
		}
	} else if node.(*Node).Children == nil {
		children, childrenErr := c.Repo.GetChildren(path)
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

func (c *CachingRepository) SaveChild(path string, child *Node) error {
	c.Invalidate(path)
	return c.Repo.SaveChild(path, child)
}

func (c *CachingRepository) SaveValue(path string, node *Node) error {
	c.Invalidate(path)
	return c.Repo.SaveValue(path, node)
}

func (c *CachingRepository) Delete(path string, node *Node) error {
	c.Invalidate(path)
	return c.Repo.Delete(path, node)
}

func (c *CachingRepository) Close() {
	c.InvalidateAll()
}

func (c *CachingRepository) Export(ctx context.Context, cancelFunc context.CancelFunc, path string) (*Node, error) {
	node, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	return node, c.doExport(ctx, cancelFunc, node, path)
}

func (c *CachingRepository) doExport(
	ctx context.Context,
	cancelFunc context.CancelFunc,
	parent *Node,
	parentPath string,
) error {
	errChan := make(chan error, 1)

	go func(ctx context.Context, cancelFunc context.CancelFunc, parent *Node, parentPath string) {
		var wg sync.WaitGroup
		for i, child := range parent.Children {
			wg.Add(1)
			go func(idx int, childName string) {
				defer wg.Done()

				childNode, err := c.Export(ctx, cancelFunc, c.Repo.buildAbsPath(parentPath, childName))
				parent.Children[idx] = childNode
				if err != nil {
					errChan <- err
					cancelFunc()
				}
			}(i, child.Name)
		}
		wg.Wait()
	}(ctx, cancelFunc, parent, parentPath)

	select {
	case <-ctx.Done():
		log.Infof("Node %s export has been canceled", parentPath)
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
