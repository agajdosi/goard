package main

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	Directories []struct {
		Dir   []string `yaml:"dir"`
		Files []string `yaml:"files"`
	}
	Files []struct {
		This string `yaml:"this"`
		That string `yaml:"that"`
	} `yaml:"files"`
}

func (c *conf) getConf(configPath string) *conf {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	fmt.Println(c.Directories)
	fmt.Println(c.Files)

	return c
}
