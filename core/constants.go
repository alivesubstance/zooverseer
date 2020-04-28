package core

const AppId = "com.github.alivesubstance.zooverseer"

// todo change to relative path
const (
	GladeFilePath      = "/home/mirian/code/zooverseer/assets/main.glade"
	ConnConfigFilePath = "/home/mirian/code/zooverseer/assets/connections.json"
	CssStyleFilePath   = "/home/mirian/code/zooverseer/assets/style.css"
)

// todo Move constants to config
// Zk Connection
const (
	ConnCacheExpireAfterAccessMinutes = 20
	ConnCacheStatsPeriodMinutes       = 5
	ConnRetryAttempts                 = 1
	ConnRetryDelay                    = 1000
)

// Zk operations
const (
	ZkCacheExpireAfterAccessMinutes = 10
	ZkCacheStatsPeriodMinutes       = 5
	ZkOpRetryAttempts               = 1
	ZkOpRetryDelay                  = 500
)

// Nodes tree and repository
const (
	NodeColumn   = 0
	NodeRootName = "/"
	NodeDummy    = "__dummy" // Dummy node to be used as real node children placeholder
)

//TODO tested and worked with relative path. Use run.sh to build and run app
//const (
//	AppId = "com.github.alivesubstance.zooverseer"
//	ConfigDir = "./config"
//)
//
//const (
//	GladeFilePath      = ConfigDir + "/main.glade"
//	ConnConfigFilePath = ConfigDir + "/connections.json"
//	ConfigFilePath = ConfigDir + "/zooverseer.yml"
//)
