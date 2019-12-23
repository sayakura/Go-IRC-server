package main

import (
	"errors"
	"fmt"
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

func ircPassHandler(client *Client, params []string, usr *User) error {
	if len(params) != 1 {
		client.send("Invalid number of parameters\n")
		return errors.New("Wrong number of parameters")
	}
	pw := params[0]
	usr.password = hashAndSalt([]byte(pw))
	return nil
}

func ircUserHandler(client *Client, params []string, usr *User) error {
	if len(params) != 1 {
		client.send("Invalid number of parameters\n")
		return errors.New("Wrong format")
	}
	username := params[0]
	usr.username = username
	return nil
}

func ircNickHandler(client *Client, params []string, usr *User) error {
	if len(params) != 1 {
		client.send("Invalid number of parameters\n")
		return errors.New("Wrong format")
	}
	nickname := params[0]
	if len(nickname) > 9 {
		client.send("nickname has a maximum length of nine (9) character\n")
		return errors.New("Wrong format")
	}
	usr.nickname = nickname
	return nil
}

func ircRegisterHandler(client *Client, params []string, usr *User) error {
	if len(params) != 3 {
		client.send("Wrong number of parameters\n")
		return errors.New("Wrong number of parameters")
	}
	usr.password = hashAndSalt([]byte(params[0]))
	if len(params[1]) > 9 {
		client.send("nickname has a maximum length of nine (9) character\n")
		return errors.New("Wrong format")
	}
	usr.nickname = params[1]
	usr.username = params[2]
	fmt.Println("good")
	return nil
}
