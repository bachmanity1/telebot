package util

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configFile = "config"

func InitConfig() *viper.Viper {
	pflag.Parse()

	var err error
	botconfig, err := readConfig(map[string]interface{}{})
	if err != nil {
		os.Exit(1)
	}

	botconfig.BindPFlags(pflag.CommandLine)
	return botconfig
}

func readConfig(defaults map[string]interface{}) (*viper.Viper, error) {
	// Read Sequence (will overloading)
	// defaults -> config file -> env -> cmd flag
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath("./util")
	v.AddConfigPath("./")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")

	v.AutomaticEnv()

	v.SetConfigName(configFile)
	return v, v.ReadInConfig()
}
