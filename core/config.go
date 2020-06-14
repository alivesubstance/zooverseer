package core

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"os"
)

const AppAssetsDir = "./assets"

type ZooverseerConfig struct {
	SortFolderFirst    bool
	ExportDir          string
	AppId              string
	UserConfigDir      string
	AppConfigDir       string
	GladeFilePath      string
	ConnConfigFilePath string
	CssStyleFilePath   string
}

var Config = ZooverseerConfig{
	AppId:              "com.github.alivesubstance.zooverseer",
	SortFolderFirst:    true,
	ConnConfigFilePath: createUserConfigDir() + "/connections.json",
	GladeFilePath:      AppAssetsDir + "/main.glade",
	CssStyleFilePath:   AppAssetsDir + "/style.css",
	ExportDir:          "./export",
}

func createUserConfigDir() string {
	userHome, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Errorf("Failed to get user home dir. Fallback to %v", AppAssetsDir)
	}

	userConfigDir := userHome + "/zooverseer"
	_, err = os.Stat(userConfigDir)
	if os.IsNotExist(err) {
		_, err = os.Create(userConfigDir)
		if err != nil {
			log.WithError(err).Fatalf("Failed to create user config dir %v", userConfigDir)
		}
	}

	return userHome
}
