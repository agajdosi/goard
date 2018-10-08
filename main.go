package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	config conf
)

func main() {
	getFlags()
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
