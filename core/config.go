package core

var Config = ZooverseerConfig{
	SortFolderFirst: true,
	ExportDir:       "/home/mirian/code/zooverseer/export",
}

type ZooverseerConfig struct {
	SortFolderFirst bool
	ExportDir       string
}
