package core

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

const AppAssetsDir = "./assets"

type ZooverseerConfig struct {
	SortFolderFirst     bool
	ExportDir           string
	AppId               string
	LogDir              string
	GladeFilePath       string
	ConnConfigFilePath  string
	CssStyleFilePath    string
	LogFilesHistorySize int
}

var Config = ZooverseerConfig{
	AppId:           "com.github.alivesubstance.zooverseer",
	SortFolderFirst: true,
	// todo not an easy task to create hidden file thru go
	//ConnConfigFilePath: getUserConfigDir() + "/connections.json",
	LogDir:              "./log",
	ConnConfigFilePath:  AppAssetsDir + "/connections.json",
	GladeFilePath:       AppAssetsDir + "/main.glade",
	CssStyleFilePath:    AppAssetsDir + "/style.css",
	ExportDir:           "./export",
	LogFilesHistorySize: 10,
}

func getUserConfigDir() string {
	userHome, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Errorf("Failed to get user home dir. Fallback to %v", AppAssetsDir)
	}

	log.Infof("User home dir is %v", userHome)

	userConfigDir := AppAssetsDir

	// todo not an easy task to create hidden file thru go
	//https://stackoverflow.com/questions/54139606/how-to-create-a-hidden-file-in-windows-mac-linux
	//userConfigDir := userHome + "/.zooverseer"
	//_, err = os.Stat(userConfigDir)
	//if os.IsNotExist(err) {
	//	err := os.Mkdir(userConfigDir, 0666)
	//	if err != nil {
	//		log.WithError(err).Fatalf("Failed to create user config dir %v", userConfigDir)
	//		userConfigDir = AppAssetsDir + "/config"
	//	}
	//}

	log.Infof("Set config dir to %v", userConfigDir)
	return userConfigDir
}
