package core

import goZk "github.com/samuel/go-zookeeper/zk"

// todo Move constants to config
// Zk Connection
const (
	ConnCacheExpireAfterAccessMinutes = 20
	ConnCacheStatsPeriodMinutes       = 5
	ConnRetryAttempts                 = 3
	ConnRetryDelay                    = 500
	ConnTimeoutSec                    = 20
)

// Zk operations
const (
	ZkCacheExpireAfterAccessMinutes = 10
	ZkCacheStatsPeriodMinutes       = 5
	ZkOpRetryAttempts               = 3
	ZkOpRetryDelay                  = 500
)

// Nodes tree and repository
const (
	NodeColumn   = 0
	NodeRootName = "/"
	NodeDummy    = "__dummy" // Dummy node to be used as real node children placeholder
)

var AclWorldAnyone = goZk.WorldACL(goZk.PermAll)
