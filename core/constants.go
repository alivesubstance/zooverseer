package core

const AppId = "com.github.alivesubstance.zooverseer"

// todo change to relative path
const (
	GladeFilePath      = "/home/mirian/code/go/src/github.com/alivesubstance/zooverseer/assets/main.glade"
	ConnConfigFilePath = "/home/mirian/code/go/src/github.com/alivesubstance/zooverseer/assets/zooverseer.json"
)

// Connection cache
const (
	ConnCacheExpireAfterAccessMinutes = 20
	ConnCacheStatsPeriodMinutes       = 5
)

// Nodes tree and repository
const (
	NodeColumn       = 0
	NodeRootTreePath = "0"
	NodeRootName     = "/"
	NodeDummy        = "__dummy" // Dummy node to be used as real node children placeholder
)
