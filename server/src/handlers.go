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

	user.password = params[0]
	if len(params[1]) > 9 {
		user.IO.send("nickname has a maximum length of nine (9) character\n")
		return errors.New("Wrong format")
	}
	user.nickname = params[1]
	user.username = params[2]
	return nil
}

func ircJoinHandler(db *DB, params []string, user *User) {
	if len(params) != 1 {
		user.IO.send("Wrong number of parameters\n")
	}

	channelName := params[0]
	if channelName[0] != '#' {
		channelName = "#" + channelName
	}
	if _, ok := db.channelList[channelName]; ok {
		c := db.channelList[channelName]
		if c.Users[user.nickname] != nil {
			user.IO.send("You are already in " + channelName + "!\n")
			return
		}
		c.Users[user.nickname] = user
		fmt.Println(len(c.Users))
		db.channelList[channelName] = c
		user.IO.send("You have joined " + channelName + "!\n")
	} else {
		ch := Channel{
			ChannelName: channelName,
			Users:       make(map[string]*User),
		}
		db.channelList[channelName] = &ch
		db.channelList[channelName].Users[user.nickname] = user
		user.IO.send("You have created " + channelName + "!\n")
	}
}

func ircListHandler(db *DB, params []string, user *User) {
	if len(params) != 0 {
		user.IO.send("Wrong number of parameters\n")
		return
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
	}
	channelName := params[0]
	if db.channelList[channelName] != nil {
		ch := db.channelList[channelName]
		if ch.Users[user.nickname] != nil {
			delete(ch.Users, user.nickname)
			if len(ch.Users) == 0 {
				delete(db.channelList, channelName)
			}
			user.IO.send("You just left the channel\n")
			return
		}
	}
	user.IO.send("Channel does not exist or you are not in that channel\n")
	return
}

func ircNamesHandler(db *DB, params []string, user *User) {
	if len(params) != 0 {
		user.IO.send("Wrong number of parameters\n")
		return
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
	}
	name := params[0]
	if name == user.nickname {
		user.IO.send("Please stop trolling\n")
		return
	}
	var msg string
	for i := 1; i < len(params); i++ {
		if i != 1 {
			msg += " "
		}
		msg += params[i]
	}

	if name[0] == '#' {
		channel := db.channelList[name]
		if channel == nil {
			user.IO.send("Channel doesnt exist!\n")
			return
		}
		if channel.Users[user.nickname] == nil {
			user.IO.send("You are not a member of this channel!\n")
			return
		}
		for _, u := range channel.Users {
			if u.nickname != user.nickname {
				fmt.Println("sendto : ", channel.ChannelName)
				u.IO.send("\n[" + channel.ChannelName + "]" + user.nickname + ": " + strings.Trim(msg, "\"") + "\n")
			}
		}
	} else { // privmsg with user
		u := db.userList[name]
		if u == nil || (u != nil && u.LoggedIn == false) {
			user.IO.send("User does not exsit / not online!\n")
		} else {
			u.IO.send("\n" + user.nickname + ": " + strings.Trim(msg, "\"") + "\n")
			return
		}
	}
}
