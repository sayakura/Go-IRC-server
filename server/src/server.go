package main 

import (
	"fmt"
	"net"
	"strings"
	"log"
	// "bufio"
	// "bytes"
	// "io"
)

const (
	CONN_HOST = "localhost"
	CONN_TYPE = "tcp"
	DELIMETER = "\n"
)

// var	userList map[string]User


func promptLoop(conn net.Conn) {
	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error when reading incomming request: ", err.Error())
	}
}

func authenticateUser(userConn net.Conn, addr string, db *DB) {
	var usr User
	buf := make([]byte, 1024)
	_, err := userConn.Read(buf)

	if err != nil {
		log.Println("Error when reading message from user: ", err.Error())
	}
	for _, msg := range strings.Split(string(buf), DELIMETER) {
		if msg[0] != 0 {
			fmt.Println("got something!")
			tokens := strings.Split(msg, " ")
			command := strings.ToUpper(tokens[0])
			params := tokens[1:]
			handler, found := authCommandList[command]
			if found {
				handler(userConn, params, &usr)
			} else {
				userConn.Write([]byte ("Unknown command"))
			}
		}
		if db.userIsMatched(usr) {
			userConn.Write([]byte ("Logged In!"))
			break 
		}
	}
}

func handleNewConnection(conn net.Conn, db *DB) {
	addr := conn.RemoteAddr().String()
	for !db.isLoggedIn(addr) {
		conn.Write([]byte("You are currently not logged in\n"))
		authenticateUser(conn, addr, db)
	}
	promptLoop(conn)
	conn.Close()
}

func runServer(cfg Config) {
	db := initDB(cfg)
	host := strings.Join([]string {CONN_HOST, cfg.Port}, ":")
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
		if err != nil {
			log.Println("Error with accepting new connection: ", err.Error())
		}
		newConn.Write([]byte ("You are connected to IRC server!\n"))
		go handleNewConnection(newConn, db)
	}
}

// var commandList = map[string]func(net.Conn, []string) {
// 	"PASS" : ircPassHandler,
// }

var authCommandList = map[string]func(net.Conn, []string, *User) {
	"PASS" : ircPassHandler,
	"NICK" : ircPassHandler,
	"USER" : ircPassHandler,
}