package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var decodeConfig viper.DecoderConfigOption = func(config *mapstructure.DecoderConfig) {
	config.TagName = "yaml"
}

func UnmarshalKey(key string, rawVal interface{}) error {
	return viper.UnmarshalKey(key, rawVal, decodeConfig)
}
