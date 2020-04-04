package zk

import (
	"github.com/alivesubstance/zooverseer/core"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	connInfo := &core.JsonConnInfo{
		Host:     "10.1.1.113",
		Port:     2181,
		User:     "zookeeper",
		Password: "z00k33p3r",
	}

	exists(connInfo, "/")
	exists(connInfo, "/env")

	log.Println("Sleep")
	time.Sleep(10 * time.Second)

	exists(connInfo, "/env/sandbox-pleeco/acl")
	exists(connInfo, "/env/sandbox-pleeco/cassandra.port")
}

func exists(connInfo *core.JsonConnInfo, path string) {
	_, stat, _ := Exists(path, connInfo)
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
