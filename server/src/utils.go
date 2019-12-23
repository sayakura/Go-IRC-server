package main

import (
	"encoding/json"
	"log"
	"os"
	//"fmt"
)

type Config struct {
	Port         string
	DatabasePath string
}

func readConfigFile(path string) Config {
	var ret Config
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("Failed to open config file")
	}
	err = json.NewDecoder(file).Decode(&ret)
	if err != nil {
		log.Fatalln("Encountered error when parsing config file ", err.Error())
	}
	return ret
}

// func hashAndSalt(pwd []byte) string {
// 	//bytes, _ := bcrypt.GenerateFromPassword([]byte(pwd), 14)
// 	//return string(bytes)
// 	return string(pwd)
// }

// func checkPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }
