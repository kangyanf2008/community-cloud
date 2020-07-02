package db

import (
	"community-cloud/config"
)

var DbInstences map[string]*GroupHandler

func init() {
	DbInstences = make(map[string]*GroupHandler, 1)
}

func GetDbGroupHandler(name string) *GroupHandler {
	if gh, ok := DbInstences[name]; ok {
		return gh
	}

	return nil
}

func InitDbConn(conf config.BaseConfig) error {
	//init db
	InitAllDB(conf, DbInstences)

	return nil
}
