package main

import (
	"fmt"
	"os"
)

// var configPath = flag.String("c", "settings.conf", "path to the configuration file")

// ./client nickname username password
func main() {
	if len(os.Args) != 5 {
		fmt.Println("usage: [host] [password] [nickname] [username]")
		return
	}
	argsWithoutProg := os.Args[1:]
	for _, a := range argsWithoutProg {
		fmt.Println(a)
	}
	host := argsWithoutProg[0]
	password := argsWithoutProg[1] 
	nickname := argsWithoutProg[2] 
	username := argsWithoutProg[3] 

	conn, _ := net.Dial("tcp", host)
	fmt.Fprintf(conn, "REGISTER %s %s %s\n", )

)

