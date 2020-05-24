package zk

import (
	"fmt"
	"github.com/alivesubstance/zooverseer/core"
	"github.com/avast/retry-go"
	goCache "github.com/goburrow/cache"
	"github.com/pkg/errors"
	goZk "github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"time"
)

var ConnCache goCache.LoadingCache

type infoLogger struct{}

func (_ infoLogger) Printf(format string, a ...interface{}) {
	log.Infof(format, a...)
}

func init() {
	connRemoveListener := func(key goCache.Key, value goCache.Value) {
		log.Debugf("Conn closed from remove listener. %s", key)
		value.(*goZk.Conn).Close()
	}

	stats := &goCache.Stats{}
	c := goCache.NewLoadingCache(connCreateFunc,
		goCache.WithExpireAfterAccess(core.ConnCacheExpireAfterAccessMinutes*time.Minute),
		goCache.WithRemovalListener(connRemoveListener),
	)
	// todo looks like stats doesn't collect numbers
	c.Stats(stats)

	go func() {
		for {
			time.Sleep(core.ConnCacheStatsPeriodMinutes * time.Minute)
			log.Infof("Conn cache: %+v\n", stats)
		}
	}()

	ConnCache = c
}

func connect(connInfo core.ConnInfo) (*goZk.Conn, error) {
	log.Infof("Connecting to %s", connInfo.String())
	goZk.DefaultLogger = &infoLogger{}

	servers := getServers(connInfo)
	conn, _, err := goZk.Connect(servers, core.ConnTimeoutSec*time.Second)
	if err != nil {
		return nil, err
	}
	err = validateConn(conn, servers)
	if err != nil {
		return nil, err
	}

	if len(connInfo.User) != 0 && len(connInfo.Password) != 0 {
		authExp := fmt.Sprint(connInfo.User, ":", connInfo.Password)
		err := conn.AddAuth("digest", []byte(authExp))
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to add auth for user %s", connInfo.User)
		}
	}

	return conn, err
}

func validateConn(conn *goZk.Conn, servers []string) error {
	data, _, err := conn.Get(core.NodeRootName)
	if data != nil && err == nil {
		return nil
	}

	// invalid connection should be closed. there is inner go routine
	// that try reconnect indefinitely. it happens f.i. in case of failed dns resolving.
	// see github.com/samuel/go-zookeeper/zk/conn.go:361
	conn.Close()
	return errors.Wrapf(err, "Failed to connect to %s", servers)
}

func getServers(connInfo core.ConnInfo) []string {
	servers := make([]string, 1)
	servers[0] = fmt.Sprintf("%v:%v", connInfo.Host, connInfo.Port)
	return servers
}

func connCreateFunc(key goCache.Key) (goCache.Value, error) {
	var validConn *goZk.Conn
	connInfo := key.(core.ConnInfo)
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
			log.WithError(err).Infof("Retry %v of %v failed to connect to %v", n, core.ConnRetryAttempts, connInfo)
		}),
	)

	return validConn, err
}
