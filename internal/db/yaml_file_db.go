package db

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/qiuhaohao/pfolio/internal/config"
)

type YamlFileDB struct {
	filePath string
	states   States
}

func NewYamlFileDB(filePath string) Database {
	return &YamlFileDB{filePath: filePath}
}

func (db *YamlFileDB) States() States {
	return db.states
}

func (db *YamlFileDB) Load() (err error) {
	db.states = newInitialStates()

	var f *os.File
	if f, err = os.Open(db.filePath); err != nil {
		return
	}

	if err = yaml.NewDecoder(f).Decode(db.states); err != nil {
		return
	}

	if err = f.Close(); err != nil {
		return
	}
	return
}

func (db *YamlFileDB) Save() (err error) {
	var f *os.File
	if f, err = os.CreateTemp("", ""); err != nil {
		return
	}
	if err = yaml.NewEncoder(f).Encode(db.states); err != nil {
		return
	}
	if err = f.Close(); err != nil {
		return
	}
	if err = os.Rename(f.Name(), viper.GetString(config.KeyDB)); err != nil {
		return
	}
	return
}
