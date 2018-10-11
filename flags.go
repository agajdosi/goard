package main

import "flag"

var (
	configPath string
	workDir    string
)

func getFlags() {
	flag.StringVar(&configPath, "config-path", "test/config.yaml", "Path to the config")
	flag.StringVar(&workDir, "work-dir", "", "Change working directory of the command. Use current directory if left empty.")
	flag.Parse()
}
