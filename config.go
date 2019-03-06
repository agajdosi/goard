package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	Directories []struct {
		Dir   []string `yaml:"dir"`
		Files []string `yaml:"files"`
	} `yaml:"directories"`

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

	return c
}

func setWorkdir() {
	if workDir != "" {
		err := os.Chdir(workDir)
		if err != nil {
			log.Fatal(err)
		}

		wd, _ := os.Getwd()
		fmt.Printf("Working directory set to: %v\n", wd)
	}
}
