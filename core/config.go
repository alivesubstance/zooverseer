package core

import (
	log "github.com/sirupsen/logrus"
)

const AppAssetsDir = "./assets"

type LogConfig struct {
	Dir              string
	Level            log.Level
	FilesHistorySize int
}

type ZooverseerConfig struct {
	ShowFolderFirst    bool
	ExportDir          string
	AppId              string
	AppTitle           string
	GladeFilePath      string
	ConnConfigFilePath string
	CssStyleFilePath   string
	LogoFilePath       string
	Log                LogConfig
}

// Because of some unbelievable reason the value "com.github.alivesubstance.zooverseer" can't be
// used as AppId. In this case "app.Run(os.Args)" returns zero code immediately and application exits. WTF!?
var Config = ZooverseerConfig{
	AppId:              "com.github.alivesubstance.app.zooverseer",
	AppTitle:           "Zooverseer",
	ShowFolderFirst:    true,
	ConnConfigFilePath: AppAssetsDir + "/connections.json",
	GladeFilePath:      AppAssetsDir + "/main.glade",
	CssStyleFilePath:   AppAssetsDir + "/style.css",
	LogoFilePath:       AppAssetsDir + "/logo.png",
	ExportDir:          "./export",
	Log: LogConfig{
		Dir:              "./log",
		Level:            log.TraceLevel,
		FilesHistorySize: 5,
	},
}
