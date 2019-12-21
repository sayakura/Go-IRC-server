package main

import (
	"encoding/json"
	"os"
	//"fmt"
)

type Config struct {
	Port string
}

func fatal(s string) {
	os.Stderr.WriteString(s)
	os.Exit(1)
}

func readConfigFile(path string) Config{
	var ret Config
	file, err := os.Open(path)
	if err != nil {
		fatal("Failed to open config file\n")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ret)
	if err != nil {
		fatal("Encountered error when parsing config file\n")
	}
	return ret
}