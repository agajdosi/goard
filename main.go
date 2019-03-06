package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

var (
	config   conf
	problems bytes.Buffer
)

func main() {
	getFlags()
	setWorkdir()
	config.getConf(configPath)

	if config.Files != nil {
		fmt.Print("Checking individualy defined files: ")
		for _, file := range config.Files {
			checkFiles(file.This, file.That)
		}
		fmt.Println(" DONE")
	}

	if config.Directories != nil {
		fmt.Print("Checking files in directories: ")
		for _, directory := range config.Directories {
			fmt.Printf("\n - %v: ", directory.Dir[0])

			if directory.Files == nil {
				var err error
				contents, err := ioutil.ReadDir(directory.Dir[0])
				if err != nil {
					log.Fatal(err)
				}

				for _, content := range contents {
					if content.IsDir() {
						continue
					}
					directory.Files = append(directory.Files, content.Name())
				}
			}

			for _, file := range directory.Files {
				this := joinPaths(directory.Dir[0], file)
				that := joinPaths(directory.Dir[1], file)
				checkFiles(this, that)
			}
		}
		fmt.Println(" DONE")
	}

	fmt.Println("\nAll checks finished")
	if problems.Len() > 0 {
		fmt.Println("FAIL - problems detected:")
		fmt.Print(problems.String())
		os.Exit(1)
	}

	fmt.Println("OK - all files up-to-date")
	os.Exit(0)
}

func checkFiles(pathOne, pathTwo string) {
	first, err := getFile(pathOne)
	if err != nil {
		logResult(err)
		return
	}

	second, err := getFile(pathTwo)
	if err != nil {
		logResult(err)
		return
	}

	if first != second {
		err := fmt.Errorf("mismatch, file '%v' is not the same as file '%v'", pathOne, pathTwo)
		logResult(err)
		return
	}

	logResult(nil)
	return
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
		err = fmt.Errorf("scheme '%v' of file '%v' not recognized", parsedURL.Scheme, parsedURL)
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
		return "", fmt.Errorf("download error, status code: '%v' when downloading file '%v'", resp.StatusCode, url)
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

func logResult(err error) {
	if err != nil {
		problems.WriteString(fmt.Sprintln(err))
		fmt.Print("x")
		return
	}

	fmt.Print(".")
	return
}
