package main

import (
	"encoding/json"
	"os"
	"log"
	//"fmt"
)

type Config struct {
	Port string
	filePath string
}

func readConfigFile(path string) Config{
	var ret Config
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("Failed to open config file")
	}
	err = json.NewDecoder(file).Decode(&ret)
	if err != nil {
		log.Fatalln("Encountered error when parsing config file")
	}
	return ret
}