package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "test/config.yaml", "Path to the config")
	flag.Parse()

	var config conf
	config.getConf(configPath)

	for _, file := range config.Files {
		same, err := checkFiles(file.This, file.That)
		if err != nil {
			fmt.Println("error occured:", err)
			continue
		}
		fmt.Println(same)
	}
}

type conf struct {
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

func checkFiles(pathOne, pathTwo string) (bool, error) {
	first, err := getFile(pathOne)
	if err != nil {
		return false, err
	}

	second, err := getFile(pathTwo)
	if err != nil {
		return false, err
	}

	equals := first == second
	return equals, nil
}

func getFile(path string) (string, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	var content string
	switch parsedURL.Scheme {
	case "http", "https":
		content, err = downloadFile(path)
	case "", "file":
		content, err = readFile(path)
	default:
		err = fmt.Errorf("unable to ")
	}

	return content, err
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting the file")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func readFile(path string) (string, error) {
	remoteFile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(remoteFile), nil
}
