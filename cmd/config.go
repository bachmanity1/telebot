package main

import (
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configFile = "config"

var telebot *viper.Viper

func init() {
	pflag.Parse()

	var err error
	telebot, err = readConfig(map[string]interface{}{})
	if err != nil {
		log.Panicf("Error when reading config: %v\n", err)
		os.Exit(1)
	}

	telebot.BindPFlags(pflag.CommandLine)
}

func readConfig(defaults map[string]interface{}) (*viper.Viper, error) {
	// Read Sequence (will overloading)
	// defaults -> config file -> env -> cmd flag
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.AddConfigPath("./")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")

	v.AutomaticEnv()

	v.SetConfigName(configFile)
	return v, v.ReadInConfig()
}
