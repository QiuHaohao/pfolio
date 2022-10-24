package config

import (
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var decodeConfig viper.DecoderConfigOption = func(config *mapstructure.DecoderConfig) {
	config.TagName = "yaml"
}

func UnmarshalKey(key string, rawVal interface{}) error {
	return viper.UnmarshalKey(key, rawVal, decodeConfig)
}

func GetKey[T any](key string) (T, error) {
	var val T
	err := UnmarshalKey(key, &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func MustGetKey[T any](key string) T {
	val, err := GetKey[T](key)
	if err != nil {
		log.Fatal(err)
	}
	return val
}
