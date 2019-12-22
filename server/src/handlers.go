package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ircPassHandler(client *Client, params []string, usr *User) {
	if len(params) != 1 {
		client.send("Invalid number of parameters\n")
	}
	pw := params[0]
	usr.password = hashAndSalt([]byte(pw))
}

func ircUserHandler(client *Client, params []string, usr *User) {
	if len(params) != 1 {
		client.send("Invalid number of parameters\n")
	}
	username := params[0]
	usr.username = username
}
