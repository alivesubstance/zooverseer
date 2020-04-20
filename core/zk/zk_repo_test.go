package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var connInfo = &core.JsonConnInfo{
	Host:     "10.1.1.112",
	Port:     2181,
	User:     "zookeeper",
	Password: "z00k33p3r",
}
var ZkRepo = Repository{}

func TestGet(t *testing.T) {
	node, _ := ZkRepo.Get("/", connInfo)
	log.Infof("%v+", node)
}

func TestExists(t *testing.T) {

	exists(connInfo, "/")
	exists(connInfo, "/env")

	log.Println("Sleep")
	time.Sleep(10 * time.Second)

	exists(connInfo, "/env/sandbox-pleeco/acl")
	exists(connInfo, "/env/sandbox-pleeco/cassandra.port")
}

func TestGetValue(t *testing.T) {
	pathToValue := make(map[string]string)
	pathToValue["/env/sandbox-pleeco/cassandra.url"] = "10.1.1.112:9042"
	pathToValue["/env/sandbox-pleeco/cassandra.port"] = "9042"
	pathToValue["/env/sandbox-pleeco/cassandra.storage.port"] = "7000"
	pathToValue["/env/sandbox-pleeco/cassandra.rpc.port"] = "9160"
	pathToValue["/env/sandbox-pleeco/cassandra.keyspace"] = "pleeco"
	pathToValue["/env/sandbox-pleeco/cassandra.dc"] = "dc1"

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		for path, expectedValue := range pathToValue {
			wg.Add(1)
			go func(path string, expectedValue string) {
				defer wg.Done()

				node, err := ZkRepo.GetValue(path, connInfo)
				if err != nil {
					log.WithError(err).Errorf("Failed to read %s", path)
				}
				assert.Equal(t, expectedValue, node.Value)
			}(path, expectedValue)
		}
	}

	log.Info("Waiting for all go routines")
	wg.Wait()
}

func TestGetChildren(t *testing.T) {
	children, _ := ZkRepo.GetChildren("/env/sandbox-pleeco", connInfo)
	for _, child := range children {
		fmt.Printf("(%v)%s\n", child.Meta.NumChildren > 0, child.Name)
	}
}

func exists(connInfo *core.JsonConnInfo, path string) {
	stat, _ := ZkRepo.GetMeta(path, connInfo)
	log.Infof("Path %s has children: %v", path, stat.NumChildren)
}

func TestSortNodesByChildrenAndName(t *testing.T) {
	core.Config.SortFolderFirst = true

	var nodes = []*Node{
		{"with-child2", "", &zkGo.Stat{NumChildren: 2}, nil},
		{"with-child1", "", &zkGo.Stat{NumChildren: 3}, nil},
		{"with-child3", "", &zkGo.Stat{NumChildren: 1}, nil},
		{"name2", "", &zkGo.Stat{NumChildren: 0}, nil},
		{"name1", "", &zkGo.Stat{NumChildren: 0}, nil},
	}

	sortNodes(nodes)
	assert.Equal(t, "with-child1", nodes[0].Name)
	assert.Equal(t, "with-child2", nodes[1].Name)
	assert.Equal(t, "with-child3", nodes[2].Name)
	assert.Equal(t, "name1", nodes[3].Name)
	assert.Equal(t, "name2", nodes[4].Name)
}

func TestSortNodesByNameOnly(t *testing.T) {
	core.Config.SortFolderFirst = false

	var nodes = []*Node{
		{"with-child2", "", &zkGo.Stat{NumChildren: 2}, nil},
		{"with-child1", "", &zkGo.Stat{NumChildren: 3}, nil},
		{"with-child3", "", &zkGo.Stat{NumChildren: 1}, nil},
		{"name2", "", &zkGo.Stat{NumChildren: 0}, nil},
		{"name1", "", &zkGo.Stat{NumChildren: 0}, nil},
	}

	sortNodes(nodes)
	assert.Equal(t, "name1", nodes[0].Name)
	assert.Equal(t, "name2", nodes[1].Name)
	assert.Equal(t, "with-child1", nodes[2].Name)
	assert.Equal(t, "with-child2", nodes[3].Name)
	assert.Equal(t, "with-child3", nodes[4].Name)
}
