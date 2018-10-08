package main

import "flag"

var (
	configPath string
)

func getFlags() {
	flag.StringVar(&configPath, "config-path", "test/config.yaml", "Path to the config")

	flag.Parse()
}
