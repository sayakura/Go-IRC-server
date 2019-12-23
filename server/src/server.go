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

func (c *Client) getAddr() string {
	return c.conn.RemoteAddr().String()
}

// var	userList map[string]User

func promptLoop(client *Client) {
	// buf := make([]byte, 1024)

	// _, err := conn.Read(buf)
	// if err != nil {
	// 	log.Println("Error when reading incomming request: ", err.Error())
	// }
}

func authenticateUser(client *Client, addr string, db *DB) error {
	var usr User
	usr.addrInfo = addr
	buf := make([]byte, 1024)
	var handlerErr error

	for {
		_, err := client.recv(buf)
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			log.Println("Error when reading message from user: ", err.Error())
		}
		for _, msg := range strings.Split(string(buf), DELIMETER) {
			if msg != "" && msg[0] != EOF {
				handlerErr = nil
				if *debug {
					fmt.Printf("Got something: [%s][%d]\n", msg, len(msg))
				}
				tokens := strings.Split(msg, " ")
				command := strings.ToUpper(tokens[0])
				params := tokens[1:]
				handler, found := authCommandList[command]
				if found {
					handlerErr = handler(client, params, &usr)
				} else {
					client.send("Unknown command\n")
				}
				if command == "REGISTER" && handlerErr != nil {
					client.send("Successfully signed up!\n")
					db.addUser(usr)
					return nil
				}
				if usr.password != "" && usr.nickname != "" && usr.username != "" && handlerErr == nil {
					if db.userIsMatched(usr) {
						client.send("Successfully logged in!\n")
						usr.LoggedIn = true
						db.addUser(usr)
					} else {
						client.send("Nickname / username / password doesn't match with the record\n")
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

func handleNewConnection(client *Client, db *DB) {
	addr := client.getAddr()
	client.send("You are currently not logged in\n")
	for !db.isLoggedIn(addr) {
		err := authenticateUser(client, addr, db)
		if err != nil {
			fmt.Println("Client disconnected")
			client.close()
			return
		}
	}
	promptLoop(client)
	client.close()
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
		var client Client = Client{conn: newConn}

		if err != nil {
			log.Println("Error with accepting new connection: ", err.Error())
		}
		newConn.Write([]byte("You are connected to IRC server!\n"))
		go handleNewConnection(&client, db)
	}
}

// var commandList = map[string]func(net.Conn, []string) {
// 	"PASS" : ircPassHandler,
// }

var authCommandList = map[string]func(*Client, []string, *User) error{
	"REGISTER": ircRegisterHandler,
	"PASS":     ircPassHandler,
	"NICK":     ircNickHandler,
	"USER":     ircUserHandler,
}
