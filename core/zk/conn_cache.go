package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/goburrow/cache"
	zkGo "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"time"
)

var ConnCache cache.LoadingCache

type infoLogger struct{}

func (_ infoLogger) Printf(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func InitConnCache() {
	connCreateFunc := func(key cache.Key) (cache.Value, error) {
		return connect(key.(core.JsonConnInfo))
	}

	connRemoveListener := func(key cache.Key, value cache.Value) {
		log.Debugf("Conn closed from remove listener. %s", key)
		value.(*zkGo.Conn).Close()
	}

	stats := cache.Stats{}
	c := cache.NewLoadingCache(connCreateFunc,
		cache.WithExpireAfterAccess(core.ConnCacheExpireAfterAccessMinutes*time.Minute),
		cache.WithRemovalListener(connRemoveListener),
	)
	c.Stats(&stats)

	go func() {
		time.Sleep(core.ConnCacheStatsPeriodMinutes * time.Minute)
		log.Infof("Conn cache: %+v\n", stats)
	}()

	ConnCache = c
}

func connect(connInfo core.JsonConnInfo) (*zkGo.Conn, error) {
	log.Infof("Connecting to %v", connInfo)
	zkGo.DefaultLogger = &infoLogger{}

	servers := getServers(connInfo)
	conn, _, err := zkGo.Connect(servers, time.Second)
	util.CheckErrorWithMsg(fmt.Sprintf("Failed to connect to %s\n", servers), err)

	if len(connInfo.User) != 0 && len(connInfo.Password) != 0 {
		authExp := fmt.Sprint(connInfo.User, ":", connInfo.Password)
		conn.AddAuth("digest", []byte(authExp))
	}

	return conn, err
}

func getServers(connInfo core.JsonConnInfo) []string {
	servers := make([]string, 1)
	servers[0] = fmt.Sprintf("%v:%v", connInfo.Host, connInfo.Port)
	return servers
}
