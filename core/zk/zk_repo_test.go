package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

var connInfo = &core.JsonConnInfo{
	Host:     "10.1.1.112",
	Port:     2181,
	User:     "zookeeper",
	Password: "z00k33p3r",
}

func TestGet(t *testing.T) {
	node, _ := Get("/", connInfo)
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

func TestGetChildren(t *testing.T) {
	children, _ := GetChildren("/env/sandbox-pleeco", connInfo)
	for _, child := range children {
		fmt.Printf("(%v)%s\n", child.Meta.NumChildren > 0, child.Name)
	}
}

func exists(connInfo *core.JsonConnInfo, path string) {
	stat, _ := GetMeta(path, connInfo)
	log.Infof("Path %s has children: %v", path, stat.NumChildren)
}

func TestSleep(t *testing.T) {
	log.Print("start")

	for i := 0; i < 10; i++ {
		log.Print(i)
		time.Sleep(300 * time.Millisecond)

		if i == 2 {
			go sleep()
		}
	}

	log.Print("end")
}

func sleep() {
	log.Print("before sleep")
	time.Sleep(1 * time.Second)
	log.Print("after sleep")
}
