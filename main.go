package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

var (
	config conf
)

func main() {
	getFlags()
	setWorkdir()
	config.getConf(configPath)

	fmt.Println("Checking individualy defined files")
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

	fmt.Println("Checking files at remote places")
	for _, dir := range config.Locations {
		for _, file := range dir.Files {
			this := joinPaths(dir.Location[0], file)
			that := joinPaths(dir.Location[1], file)
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

	fmt.Println("Checking files in local directories")
	for _, directory := range config.Directories {
		contents, err := ioutil.ReadDir(directory.Dir[0])
		if err != nil {
			log.Fatal(err)
		}
		for _, content := range contents {
			if content.IsDir() {
				continue
			}
			this := joinPaths(directory.Dir[0], content.Name())
			that := joinPaths(directory.Dir[1], content.Name())
			same, err := checkFiles(this, that)
			if err != nil {
				log.Fatalf("error checking the files '%v', '%v': %v", this, that, err)
				continue
			}
			if same != true {
				fmt.Printf("mismatch: '%v', '%v'\n", this, that)
			}
		}
	}
	fmt.Println()

	fmt.Println("All checks finished")
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
