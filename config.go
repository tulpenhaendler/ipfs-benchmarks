package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Aws   Aws   `yaml:"aws"`
	Nodes Nodes `yaml:"nodes"`
}

type Aws struct {
	Secret interface{} `yaml:"secret"`
	Key    interface{} `yaml:"key"`
}

type Nodes struct {
	Regions  []Instances `yaml:"instances"`
	Instance string      `yaml:"instance"`
}

type Instances struct {
	Count  int    `yaml:"count"`
	Region string `yaml:"region"`
}

func NewConfig() *Config {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/ipfs-bench/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.ipfs-bench")  // call multiple times to add many search paths
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	c := Config{}
	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Println("unable to decode into struct, %v", err)
		os.Exit(1)
	}
	return &c
}