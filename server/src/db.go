package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Channel struct {
	ChannelName string
	Users       map[string]*User
}

type User struct {
	IO       *Client
	addrInfo string
	LoggedIn bool
	username string
	nickname string
	password string
}

type DB struct {
	userList    map[string]*User
	channelList map[string]*Channel
}

func initDB(cfg Config) *DB {
	db := new(DB)
	if *dataPresist {
		file, err := os.Open(cfg.DatabasePath)
		if err != nil {
			log.Fatalln("Failed to open data file", err.Error())
		}
		err = json.NewDecoder(file).Decode(&db)
		if err != nil && err != io.EOF {
			log.Fatalln("Encountered error when parsing data file", err.Error())
		}
	}
	if db.userList == nil {
		db.userList = make(map[string]*User)
	}
	if db.channelList == nil {
		db.channelList = make(map[string]*Channel)
	}
	return db
}

func (d *DB) isLoggedIn(user *User) bool {
	//fmt.Println(d.userList)
	for _, u := range d.userList {
		if u.nickname == user.nickname &&
			u.LoggedIn {
			return true
		}
	}
	return false
}

func (d *DB) addUser(user *User) {
	// user.password = user.password
	d.userList[user.nickname] = user
}

// }

func (d *DB) login(user *User) {
	for _, u := range d.userList {
		if u.nickname == user.nickname {
			u.LoggedIn = true
			u.IO = user.IO
		}
	}
}

func (d *DB) ifNicknameTaken(nickname string) bool {
	for _, u := range d.userList {
		if u.nickname == nickname {
			return true
		}
	}
	return false
}

func (d *DB) userIsMatched(curUser *User) bool {
	for _, u := range d.userList {
		//fmt.Println(u)
		if u.nickname == curUser.nickname &&
			u.password == curUser.password &&
			u.username == curUser.username {
			return true
		}
	}
	return false
}

func (d *DB) savedToFileSystem() {
	fmt.Println(d)
	file, _ := json.Marshal(d)
	fmt.Println(string(file))
	err := ioutil.WriteFile("./data/db", file, 0644)
	if err != nil {
		log.Fatalln("Failed to save db data info file system", err.Error())
	}
}
