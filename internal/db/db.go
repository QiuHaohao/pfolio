package db

import (
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/qiuhaohao/pfolio/internal/config"
)

var defaultDB Database

func DefaultDB() Database {
	return defaultDB
}

type Database interface {
	States() States
	Load() error
	Save() error
}

func LoadDefaultDB() {
	dbFilePath := viper.GetString(config.KeyDB)
	defaultDB = NewYamlFileDB(dbFilePath)
	if err := defaultDB.Load(); os.IsNotExist(err) {
		log.Printf("DB file %s not found, initializing with new initial states", dbFilePath)
	} else if err != nil {
		log.Fatal(err)
	}
}
