package core

import zookeeper "github.com/outbrain/zookeepercli/go/zk"

var zk = zookeeper.NewZooKeeper()

//logger, _ := zap.NewProduction()

func initZk(info JsonConnInfo) {

}

func Get(path string, info JsonConnInfo) {
	//zk.Get()
	//zook.SetServers()
	//zook.Get("/")
}
