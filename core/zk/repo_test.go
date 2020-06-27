package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// todo consider to rewrite test with https://github.com/stretchr/testify

// localhost
var connInfo = &core.ConnInfo{Host: "127.0.0.1", Port: 2181}

// sandbox-pleeco
//var connInfo = &core.ConnInfo{Host: "10.1.1.112", Port: 2181, User: "zookeeper", Password: "z00k33p3r"}

var ZkRepo = Repository{connInfo: connInfo}

func TestGet(t *testing.T) {
	node, _ := ZkRepo.Get("/")
	log.Infof("%v+", node)
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

				node, err := ZkRepo.GetValue(path)
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
	children, _ := ZkRepo.GetChildren("/env/sandbox-pleeco")
	for _, child := range children {
		fmt.Printf("(%v)%s\n", child.Meta.NumChildren > 0, child.Name)
	}
}

func TestSave(t *testing.T) {
	ZkRepo.SetConnInfo(connInfo)

	nodeName := "test"
	err := ZkRepo.SaveChild("/", nodeName, core.AclWorldAnyone)
	assert.Nil(t, err)
	if err != nil {
		log.WithError(err).Panicf("Failed to save node " + nodeName)
	}

	meta, _ := ZkRepo.GetMeta("/" + nodeName)
	assert.NotNil(t, meta)
}

func TestSortNodesByChildrenAndName(t *testing.T) {
	core.Config.ShowFolderFirst = true

	var nodes = []*Node{
		{"with-child2", "", &zkGo.Stat{NumChildren: 2}, nil, nil},
		{"with-child1", "", &zkGo.Stat{NumChildren: 3}, nil, nil},
		{"with-child3", "", &zkGo.Stat{NumChildren: 1}, nil, nil},
		{"name2", "", &zkGo.Stat{NumChildren: 0}, nil, nil},
		{"name1", "", &zkGo.Stat{NumChildren: 0}, nil, nil},
	}

	sortNodes(nodes)
	assert.Equal(t, "with-child1", nodes[0].Name)
	assert.Equal(t, "with-child2", nodes[1].Name)
	assert.Equal(t, "with-child3", nodes[2].Name)
	assert.Equal(t, "name1", nodes[3].Name)
	assert.Equal(t, "name2", nodes[4].Name)
}

func TestSortNodesByNameOnly(t *testing.T) {
	core.Config.ShowFolderFirst = false

	var nodes = []*Node{
		{"with-child2", "", &zkGo.Stat{NumChildren: 2}, nil, nil},
		{"with-child1", "", &zkGo.Stat{NumChildren: 3}, nil, nil},
		{"with-child3", "", &zkGo.Stat{NumChildren: 1}, nil, nil},
		{"name2", "", &zkGo.Stat{NumChildren: 0}, nil, nil},
		{"name1", "", &zkGo.Stat{NumChildren: 0}, nil, nil},
	}

	sortNodes(nodes)
	assert.Equal(t, "name1", nodes[0].Name)
	assert.Equal(t, "name2", nodes[1].Name)
	assert.Equal(t, "with-child1", nodes[2].Name)
	assert.Equal(t, "with-child2", nodes[3].Name)
	assert.Equal(t, "with-child3", nodes[4].Name)
}
