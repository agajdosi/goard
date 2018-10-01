package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	result, _ := checkFiles("test/eap71-basic-s2i.json", "https://raw.githubusercontent.com/minishift/minishift/master/addons/xpaas/v3.10/xpaas-templates/eap71-basic-s2i.json")
	fmt.Println(result)

	result, _ = checkFiles("test/eap71-basic-s2i-wrong.json", "https://raw.githubusercontent.com/minishift/minishift/master/addons/xpaas/v3.10/xpaas-templates/eap71-basic-s2i.json")
	fmt.Println(result)
}

func checkFiles(path, url string) (bool, error) {
	remoteFile, err := getFile(url)
	if err != nil {
		return false, err
	}

	localFile, err := readFile(path)
	if err != nil {
		return false, err
	}

	result := areSame(localFile, remoteFile)

	return result, nil
}

func getFile(url string) (string, error) {
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

func areSame(stringOne, stringTwo string) bool {
	if stringOne != stringTwo {
		return false
	}

	return true
}
