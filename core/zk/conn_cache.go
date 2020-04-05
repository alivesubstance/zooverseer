package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/alivesubstance/zooverseer/util"
	"github.com/avast/retry-go"
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

func connCreateFunc(key cache.Key) (cache.Value, error) {
	var validConn *zkGo.Conn
	connInfo := key.(core.JsonConnInfo)
	err := retry.Do(
		func() error {
			conn, err := connect(connInfo)
			if err != nil {
				return err
			}

			validConn = conn
			return nil
		},
		retry.Attempts(core.ConnRetryAttempts),
		retry.Delay(core.ConnRetryDelay*time.Millisecond),
		retry.OnRetry(func(n uint, err error) {
			log.WithError(err).Infof("Retry %v of %v failed to connect to %s", n, core.ConnRetryAttempts, connInfo)
		}),
	)

	return validConn, err
}
