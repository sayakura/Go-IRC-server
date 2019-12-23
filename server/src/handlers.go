package main

import (
	"errors"
	"fmt"
	"strings"
)

func ircPassHandler(params []string, user *User) error {
	if len(params) != 1 {
		user.IO.send("Invalid number of parameters\n")
		return errors.New("Wrong number of parameters")
	}
	pw := params[0]
	user.password = pw
	return nil
}

func ircUserHandler(params []string, user *User) error {
	if len(params) != 1 {
		user.IO.send("Invalid number of parameters\n")
		return errors.New("Wrong format")
	}
	username := params[0]
	user.username = username
	return nil
}

func ircNickHandler(params []string, user *User) error {
	if len(params) != 1 {
		user.IO.send("Invalid number of parameters\n")
		return errors.New("Wrong format")
	}
	nickname := params[0]
	if len(nickname) > 9 {
		user.IO.send("nickname has a maximum length of nine (9) character\n")
		return errors.New("Wrong format")
	}
	user.nickname = nickname
	return nil
}

func ircRegisterHandler(params []string, user *User) error {
	if len(params) != 3 {
		user.IO.send("Wrong number of parameters\n")
		return errors.New("Wrong number of parameters")
	}

	user.password = hashAndSalt([]byte(params[0]))
	if len(params[1]) > 9 {
		user.IO.send("nickname has a maximum length of nine (9) character\n")
		return errors.New("Wrong format")
	}
	user.nickname = params[1]
	user.username = params[2]
	fmt.Println("good")
	return nil
}

func ircJoinHandler(db *DB, params []string, user *User) {
	if len(params) != 1 {
		user.IO.send("Wrong number of parameters\n")
		//return errors.New("Wrong number of parameters")
	}

	channelName := params[0]
	if channelName[0] != '#' {
		channelName = "#" + channelName
	}
	if _, ok := db.channelList[channelName]; ok {
		c := db.channelList[channelName]
		for _, u := range c.Users {
			if u.nickname == user.nickname {
				user.IO.send("You are already in " + channelName + "!\n")
				return
			}
		}
		c.Users = append(c.Users, user)
		fmt.Println(len(c.Users))
		db.channelList[channelName] = c
		user.IO.send("You have joined " + channelName + "!\n")
	} else {
		db.channelList[channelName] = Channel{
			ChannelName: channelName,
			Users:       make([]*User, 1),
		}
		db.channelList[channelName].Users[0] = user
		user.IO.send("You have created " + channelName + "!\n")
	}
}

func ircListHandler(db *DB, params []string, user *User) {
	if len(params) != 0 {
		user.IO.send("Wrong number of parameters\n")
		return
		//return errors.New("Wrong number of parameters")
	}
	var s string
	for _, v := range db.channelList {
		s += v.ChannelName + "\n"
	}
	if s == "" {
		user.IO.send("No channel available\n")
	} else {
		user.IO.send(s)
	}
}

func ircPartHandler(db *DB, params []string, user *User) {
	if len(params) != 1 {
		user.IO.send("Wrong number of parameters\n")
		return
		//return errors.New("Wrong number of parameters")
	}
	channelName := params[0]
	for _, v := range db.channelList {
		if channelName == v.ChannelName {
			for i, u := range v.Users {
				if u.nickname == user.nickname {
					v.Users = append(v.Users[:i], v.Users[i+1:]...)
					if len(v.Users) == 0 {
						delete(db.channelList, channelName)
					}
					user.IO.send("You have leave the channel\n")
					return
				}
			}
		}
	}
	user.IO.send("Channel does not exist or you are not in that channel\n")
	return
}

func ircNamesHandler(db *DB, params []string, user *User) {
	if len(params) != 0 {
		user.IO.send("Wrong number of parameters\n")
		return
		//return errors.New("Wrong number of parameters")
	}
	var s string
	for _, v := range db.userList {
		if v.LoggedIn {
			s += v.nickname + "\n"
		}
	}
	if s == "" {
		user.IO.send("No users available\n")
	} else {
		user.IO.send(s)
	}
}

func ircPrivmsgHandler(db *DB, params []string, user *User) {
	if len(params) < 2 {
		user.IO.send("Wrong number of parameters\n")
		return
		//return errors.New("Wrong number of parameters")
	}
	name := params[0]
	var msg string
	for i := 1; i < len(params); i++ {
		if i != 1 {
			msg += " "
		}
		msg += params[i]
	}
	if name[0] == '#' {
		var channel *Channel = nil
		for _, c := range db.channelList {
			if c.ChannelName == name {
				channel = &c
			}
		}
		if channel == nil {
			user.IO.send("Channel doesnt exist!\n")
			return
		}
		fmt.Println("sending to everyone on the channel!...")
		for _, u := range channel.Users {
			if u.nickname != user.nickname {
				u.IO.send("\n[" + channel.ChannelName + "]" + user.nickname + ": " + strings.Trim(msg, "\"") + "\n")
			}
		}
	} else {
		for _, u := range db.userList {
			if u.nickname == name {
				u.IO.send("\n" + user.nickname + ": " + strings.Trim(msg, "\"") + "\n")
				return
			}
		}
	}
}
