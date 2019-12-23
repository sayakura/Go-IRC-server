package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	// "bufio"
	// "bytes"
	"io"
)

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
	DELIMETER = "\n"
	EOF       = 0
)

type Client struct {
	conn net.Conn
}

func (c *Client) send(msg string) {
	c.conn.Write([]byte(msg))
}

func (c *Client) recv(buf []byte) (int, error) {
	return c.conn.Read(buf)
}

func (c *Client) close() {
	c.conn.Close()
}

func promptLoop(user *User, db *DB) error {
	buf := make([]byte, 1024)

	for {
		_, err := user.IO.recv(buf)
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			log.Println("Error when reading message from user: ", err.Error())
		}
		for _, msg := range strings.Split(string(buf), DELIMETER) {
			if msg != "" && msg[0] != EOF {
				msg = strings.Trim(msg, " ")
				tokens := strings.Split(msg, " ")
				command := strings.ToUpper(tokens[0])
				params := tokens[1:]
				handler, found := commandList[command]
				if found {
					handler(db, params, user)
				} else {
					user.IO.send("Unknown command\n")
				}
			}
		}
		for i := 0; i < 1024; i++ {
			buf[i] = 0
		}
	}
}

func authenticateUser(user *User, db *DB) error {
	buf := make([]byte, 1024)
	var handlerErr error

	for {
		_, err := user.IO.recv(buf)
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			log.Println("Error when reading message from user: ", err.Error())
		}
		for _, msg := range strings.Split(string(buf), DELIMETER) {
			if msg != "" && msg[0] != EOF {
				msg = strings.Trim(msg, " ")
				handlerErr = nil
				tokens := strings.Split(msg, " ")
				command := strings.ToUpper(tokens[0])
				params := tokens[1:]
				handler, found := authCommandList[command]
				if found {
					handlerErr = handler(params, user)
				} else {
					user.IO.send("Unknown command\n")
				}
				if command == "REGISTER" && handlerErr == nil {
					if db.ifNicknameTaken(user.nickname) {
						user.IO.send("Nickname taken, choose a different one\n")
						user.nickname = ""
						user.password = ""
						user.username = ""
						return nil
					}
					user.IO.send("Successfully signed up and logged in!\n")
					user.LoggedIn = true
					db.addUser(user)
					return nil
				}
				if user.password != "" && user.nickname != "" && user.username != "" && handlerErr == nil {
					if db.userIsMatched(user) {
						user.IO.send("Successfully logged in!\n")
						db.login(user)
					} else {
						user.IO.send("Nickname / username / password doesn't match with the record\n")
					}
					return nil
				}
			}
		}
		for i := 0; i < 1024; i++ {
			buf[i] = 0
		}
	}
}

func handleNewConnection(user *User, db *DB) {
	user.IO.send("You are currently not logged in\n")
	for !db.isLoggedIn(user) {
		err := authenticateUser(user, db)
		if err != nil {
			fmt.Println("Client disconnected")
			user.IO.close()
			return
		}
	}
	promptLoop(user, db)
	db.userList[user.nickname].LoggedIn = false
	user.IO.close()
}

func runServer(cfg Config, db *DB) {
	host := strings.Join([]string{CONN_HOST, cfg.Port}, ":")
	if *debug {
		fmt.Printf("Host: %s\n", host)
	}
	conn, err := net.Listen(CONN_TYPE, host)
	if err != nil {
		log.Fatal("Can't listen to port ", cfg.Port, " at localhost\n")
	}
	defer conn.Close()

	fmt.Printf("Starting server...\n")
	fmt.Printf("Listening on Port %s\n", cfg.Port)
	for {
		newConn, err := conn.Accept()
		var user User = User{
			IO:       &Client{conn: newConn},
			addrInfo: newConn.RemoteAddr().String(),
		}

		if err != nil {
			log.Println("Error with accepting new connection: ", err.Error())
		} else {
			if *debug {
				fmt.Println("New connection!: ", user.addrInfo)
			}
		}
		newConn.Write([]byte("You are connected to IRC server!\n"))
		go handleNewConnection(&user, db)
	}
}

var commandList = map[string]func(*DB, []string, *User){
	"JOIN":    ircJoinHandler,
	"PART":    ircPartHandler,
	"NAMES":   ircNamesHandler,
	"LIST":    ircListHandler,
	"PRIVMSG": ircPrivmsgHandler,
}

var authCommandList = map[string]func([]string, *User) error{
	"REGISTER": ircRegisterHandler,
	"PASS":     ircPassHandler,
	"NICK":     ircNickHandler,
	"USER":     ircUserHandler,
}
