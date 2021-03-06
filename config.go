package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

type Config struct {
	Nodes Nodes `yaml:"nodes"`
	RunCmd string `yaml:"runCmd"`
}


type Nodes struct {
	Instances    []Instances `yaml:"instances"`
	InstanceType string      `yaml:"instanceType"`
}

type Instances struct {
	Count  int    `yaml:"count"`
	Region string `yaml:"region"`
	Name string `yaml:"name"`
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

func (c *Config) GetNumRegions() int {
	r := map[string]struct{}{}
	for _,a := range c.Nodes.Instances {
		r[a.Region] = struct{}{}
	}
	return len(r)
}

func (c *Config) GetRegions() []string {
	r := map[string]struct{}{}
	for _,a := range c.Nodes.Instances {
		r[a.Region] = struct{}{}
	}
	res := []string{}
	for a,_ := range r {
		res  = append(res, a)
	}
	return res
}

func (c *Config) GetInstances() []Instances {
	res := []Instances{}
	for _,a := range c.Nodes.Instances {
		i := 1
		for {
			res = append(res, Instances{
				Name: a.Name + "_" + strconv.Itoa(i),
				Region: a.Region,
				Count: 1,
			})
			if i >= a.Count {
				break
			}
			i++
		}
	}
	return res
}

func (c *Config) getNumNodes() int {
	return len(c.Nodes.Instances)
}