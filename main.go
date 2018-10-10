package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

var (
	config conf
)

func main() {
	getFlags()
	config.getConf(configPath)

	fmt.Println("Checking individualy defined files...")
	for _, file := range config.Files {
		same, err := checkFiles(file.This, file.That)
		if err != nil {
			log.Fatalf("error checking the files '%v', '%v': %v", file.This, file.That, err)
			continue
		}
		if same != true {
			log.Fatalf("files mismatch: '%v', '%v'\n", file.This, file.That)
		}
	}

	fmt.Println("Checking files in directories...")
	for _, dir := range config.Directories {
		for _, file := range dir.Files {
			this := joinPaths(dir.Dir[0], file)
			that := joinPaths(dir.Dir[1], file)
			same, err := checkFiles(this, that)
			if err != nil {
				log.Fatalf("error checking the files '%v', '%v': %v", this, that, err)
				continue
			}
			if same != true {
				log.Fatalf("files mismatch: '%v', '%v'\n", this, that)
			}
		}
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

func joinPaths(first string, second string) string {
	var result string
	firstURL, err := url.Parse(first)
	if err != nil {
		log.Fatal(err)
	}

	switch firstURL.Scheme {
	case "http", "https":
		firstURL.Path = path.Join(firstURL.Path, second)
		result = firstURL.String()
	case "", "file":
		result = path.Join(first, second)
	default:
		log.Fatalf("unable to match the path '%v' to http(s):// or file:// schemes", first)
	}

	return result
}
