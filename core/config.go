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

var Config = ZooverseerConfig{
	AppId:              "com.github.alivesubstance.zooverseer",
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
