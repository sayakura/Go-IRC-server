package main

import (
	"encoding/json"
	"os"
	"log"
)

type Channel struct {
	
}

type User struct {
	addrInfo		string
	LoggedIn		bool
	username		string
	nickname		string
	password		string
	channels		[]Channel
}

type DB struct {
	userList map[string]User
	channelList []Channel
}

func initDB(cfg Config) *DB{
	db := new(DB)
	if *dataPresist {
		file, err := os.Open(cfg.filePath)
		if err != nil {
			log.Fatalln("Failed to open data file")
		}
		err = json.NewDecoder(file).Decode(&db)
		if err != nil {
			log.Fatalln("Encountered error when parsing data file")
		}
	}
	return db
}

func (d *DB) isLoggedIn(addr string) bool {
	return d.userList[addr].LoggedIn
}

func (d *DB) userIsMatched(curUser User) bool {
	for _, user := range d.userList {
		if (curUser.username == user.username || curUser.nickname == user.nickname) && 
			curUser.password == user.password {
			return true
		}
	}
	return false
}